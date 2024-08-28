package handlers

import (
	"ChoreQuest/internal/models"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func RewardListHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rewards, err := models.GetAllRewards(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("../../templates/reward_list.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, rewards)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func AddRewardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl, err := template.ParseFiles("../../templates/add_reward.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodPost {
			description := r.FormValue("description")
			pointCost, _ := strconv.Atoi(r.FormValue("point-cost"))

			reward := &models.Reward{
				Description: description,
				PointCost:   pointCost,
			}

			err := reward.Save(db)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("HX-Trigger", "refreshRewardList")
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<button class="action-button" hx-get="/add-reward" hx-target="#reward-action-container" hx-swap="innerHTML">Add Reward</button>`))
		}
	}
}

func EditRewardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paths := strings.Split(r.URL.Path, "/")
		id, err := strconv.ParseInt(paths[len(paths)-1], 10, 64)
		if err != nil {
			http.Error(w, "Invalid reward ID", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			reward, err := models.GetRewardByID(db, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tmpl, err := template.ParseFiles("../../templates/edit_reward.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, reward)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodPost {
			description := r.FormValue("description")
			pointCost, _ := strconv.Atoi(r.FormValue("point-cost"))

			reward := &models.Reward{
				ID:          id,
				Description: description,
				PointCost:   pointCost,
			}

			err := reward.Save(db)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("HX-Trigger", "refreshRewardList")
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<button class="action-button" hx-get="/add-reward" hx-target="#reward-action-container" hx-swap="innerHTML">Add Reward</button>`))
		}
	}
}

func DeleteRewardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		paths := strings.Split(r.URL.Path, "/")
		id, err := strconv.ParseInt(paths[len(paths)-1], 10, 64)
		if err != nil {
			http.Error(w, "Invalid reward ID", http.StatusBadRequest)
			return
		}

		err = models.DeleteReward(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func RewardActionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<button class="action-button" hx-get="/add-reward" hx-target="#reward-action-container" hx-swap="innerHTML">Add Reward</button>`))
	}
}

func RewardsStoreHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract childID from URL
		paths := strings.Split(r.URL.Path, "/")
		childIDStr := paths[len(paths)-1]

		childID, err := strconv.ParseInt(childIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid child ID", http.StatusBadRequest)
			return
		}

		child, err := models.GetChildByID(db, childID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rewards, err := models.GetAllRewards(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Sort rewards by point cost
		sort.Slice(rewards, func(i, j int) bool {
			return rewards[i].PointCost < rewards[j].PointCost
		})

		data := struct {
			Child   *models.Child
			Rewards []*models.Reward
		}{
			Child:   child,
			Rewards: rewards,
		}

		tmpl, err := template.ParseFiles("../../templates/rewards_store.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func RedeemRewardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		childID, err := strconv.ParseInt(r.FormValue("child_id"), 10, 64)
		if err != nil {
			log.Printf("Error: Invalid child ID: %v", err)
			http.Error(w, "Invalid child ID", http.StatusBadRequest)
			return
		}

		rewardID, err := strconv.ParseInt(r.FormValue("reward_id"), 10, 64)
		if err != nil {
			log.Printf("Error: Invalid reward ID: %v", err)
			http.Error(w, "Invalid reward ID", http.StatusBadRequest)
			return
		}

		child, err := models.GetChildByID(db, childID)
		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		reward, err := models.GetRewardByID(db, rewardID)
		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if child.Points < reward.PointCost {
			log.Printf("Error: Not enough points")
			http.Error(w, "Not enough points", http.StatusBadRequest)
			return
		}

		// Deduct points and add reward to child
		child.Points -= reward.PointCost
		if child.Rewards == "" {
			child.Rewards = reward.Description
		} else {
			child.Rewards += ", " + reward.Description
		}

		err = child.Save(db)
		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect back to rewards store page
		log.Printf("Redirecting to /rewards-store/%d", childID)
		http.Redirect(w, r, fmt.Sprintf("/rewards-store/%d", childID), http.StatusSeeOther)
	}
}
