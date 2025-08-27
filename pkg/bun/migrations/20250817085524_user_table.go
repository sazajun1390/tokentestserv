package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            int64     `bun:"id,pk,autoincrement"`
	CreatedAt     time.Time `bun:"created_at,notnull"`

	UserProfiles *UserProfile `bun:"rel:has-one,join:id=user_id"`
	UserActives  *UserActive  `bun:"rel:has-one,join:id=user_id"`
}

type UserStatus string

const (
	UserStatusActive       UserStatus = "active"
	UserStatusProvisioning UserStatus = "provisioning"
	UserStatusInactive     UserStatus = "inactive"
	UserStatusDeleted      UserStatus = "deleted"
	UserStatusPurged       UserStatus = "purged"
	UserStatusUnspecified  UserStatus = "unspecified"
)

type UserProfile struct {
	bun.BaseModel  `bun:"table:user_profiles,alias:up"`
	UserID         int64          `bun:"user_id,notnull,unique"`
	UserMultiID    string         `bun:"user_multi_id,notnull,unique"`
	ResourceID     string         `bun:"resource_id,notnull,unique"`
	Email          string         `bun:"email,unique,notnull"`
	Password       string         `bun:"password,notnull"`
	Tel            sql.NullString `bun:"tel"`
	CreatedAt      time.Time      `bun:"created_at,notnull"`
	UpdatedAt      time.Time      `bun:"updated_at,notnull"`
	DeletedAt      bun.NullTime   `bun:"deleted_at,soft_delete"`
	PurgeExpiredAt bun.NullTime   `bun:"purged_expires_at"`

	// User is the user that this profile belongs to.
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

type UserActive struct {
	bun.BaseModel `bun:"table:user_actives,alias:ua"`
	UserID        int64      `bun:"user_id,notnull,unique"`
	Status        UserStatus `bun:"type:status_enum,notnull,default:'active'"`
	CreatedAt     time.Time  `bun:"created_at,notnull"`
	UpdatedAt     time.Time  `bun:"updated_at,notnull"`

	// User is the user that this active belongs to.
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		var err error
		_, err = db.Exec(`
            CREATE TYPE status_enum AS ENUM ('active', 'provisioning', 'inactive', 'deleted', 'purged', 'unspecified');
        `)
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().Model((*User)(nil)).WithForeignKeys().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().Model((*UserProfile)(nil)).WithForeignKeys().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().Model((*UserActive)(nil)).WithForeignKeys().Exec(ctx)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		_, err := db.NewDropTable().Model((*UserActive)(nil)).IfExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewDropTable().Model((*UserProfile)(nil)).IfExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewDropTable().Model((*User)(nil)).IfExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.Exec(`DROP TYPE status_enum;`)
		if err != nil {
			return err
		}
		return nil
	})
}
