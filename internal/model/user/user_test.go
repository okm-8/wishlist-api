package user

import "testing"

func TestUser(t *testing.T) {
	t.Run("should promote admin", func(t *testing.T) {
		t.Parallel()

		user := New("test@test.com", "test")

		user.PromoteAdmin()

		if !user.IsAdmin() {
			t.Errorf("expected admin, got not admin")
		}
	})

	t.Run("should demote admin", func(t *testing.T) {
		t.Parallel()

		user := New("test@test.com", "test")

		user.PromoteAdmin()
		user.DemoteAdmin()

		if user.IsAdmin() {
			t.Errorf("expected not admin, got admin")
		}
	})

	t.Run("should change email", func(t *testing.T) {
		t.Parallel()

		user := New("test@test.com", "test")

		user.ChangeEmail("new@test.com")

		if user.Email() != "new@test.com" {
			t.Errorf("unexpected email: %s", user.Email())
		}
	})

	t.Run("should change name", func(t *testing.T) {
		t.Parallel()

		user := New("test@test.com", "test")

		user.ChangeName("new")

		if user.Name() != "new" {
			t.Errorf("unexpected name: %s", user.Name())
		}
	})

	t.Run("should restore user", func(t *testing.T) {
		t.Parallel()

		user := New("test@test.com", "test")

		restored := Restore(user.Id(), user.Email(), user.Name(), user.IsAdmin(), user.PasswordHash())

		if !user.Equal(restored) {
			t.Errorf("expected restored user to be equal to original user")
		}
	})
}
