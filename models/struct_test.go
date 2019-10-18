package models

import "testing"

func TestStateValues(t *testing.T) {
	if UserTypeAdministrator != 0 ||
		UserTypeNormalUser != 1 ||

		UserStateNormal != 0 ||
		UserStateLock != 1 ||

		PostStateNormal != 0 ||
		PostStateLock != 1 ||
		PostStateHide != 2 ||

		CmtStateNormal != 0 ||
		CmtStateLock != 1 ||
		CmtStateHide != 2 {

		t.Error("UserTypeAdministrator:")
		t.Error(UserTypeAdministrator)
		t.Error("UserTypeNormalUser:")
		t.Error(UserTypeNormalUser)

		t.Error("UserStateNormal:")
		t.Error(UserStateNormal)
		t.Error("UserStateLock:")
		t.Error(UserStateLock)

		t.Error("PostStateNormal:")
		t.Error(PostStateNormal)
		t.Error("PostStateLock:")
		t.Error(PostStateLock)
		t.Error("PostStateHide:")
		t.Error(PostStateHide)

		t.Error("CmtStateNormal:")
		t.Error(CmtStateNormal)
		t.Error("CmtStateLock:")
		t.Error(CmtStateLock)
		t.Error("CmtStateHide:")
		t.Error(CmtStateHide)

		t.Error("状态、类型数值错误")
		t.FailNow()
	}
}
