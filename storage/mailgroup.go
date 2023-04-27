package storage

import (
	"strings"

	"github.com/pkg/errors"
)

// CreateMaillistOfUser this func will accept maillist string array as input, delete all entry of the user in users_mailgroup, and insert each input entry into users_mailgroup table
func (gs *GormStorage) CreateMaillistOfUser(uid string, maillist []string) error {
	// delete all entry of the user in users_mailgroup
	err := gs.db.Exec("DELETE FROM users_mailgroups WHERE user_id = ?", uid).Error
	if err != nil {
		return errors.Wrap(err, "insert user mailgroup error")
	}

	// insert new entry into users_mailgroup
	for _, list := range maillist {
		err = gs.db.Exec("INSERT INTO users_mailgroups (user_id, mailgroup_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE mailgroup_id = ?", uid, list, list).Error
		if err != nil {
			return errors.Wrap(err, "insert user mailgroup error")
		}
	}

	// Get the user's email from the user table
	var userEmail []string
	err = gs.db.Raw("SELECT email FROM users WHERE id = ?", uid).Pluck("email", &userEmail).Error
	if err != nil {
		return errors.Wrap(err, "get user email error")
	}

	// Convert maillist array to comma-separated string
	maillistStr := strings.Join(maillist, ",")

	// Insert new entry into jobs_mailchimp table
	err = gs.db.Exec("INSERT INTO jobs_mailchimp (receiver, interests, state) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE interests = ?", userEmail[0], maillistStr, "new", maillistStr).Error
	if err != nil {
		return errors.Wrap(err, "insert user mailgroup error")
	}

	return nil
}
