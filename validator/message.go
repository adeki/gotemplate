package validator

var msgMap = map[string]map[string]string{
	"fields": {
		"email":    "メールアドレス",
		"name":     "名前",
		"password": "パスワード",
		"url":      "URL",
	},
	"tags": {
		"exists":   "%sはすでに登録済みです。",
		"invalid":  "%sが不正です",
		"max":      "%sの入力値が大きすぎます。",
		"min":      "%sの入力値が小さすぎます。",
		"ne":       "%sが不正です",
		"oneof":    "%sの選択肢が不正です。",
		"past":     "%sは現在日時より後の日付を入力してください",
		"required": "%sを入力して下さい。",
	},
	"messages": {
		"account.innactive": "無効なアカウントです。",
		"email.eqfield":     "%sと確認用メールアドレスが違います。",
		"name.max":          "%sは最大255文字まで入力可能です。",
		"password.eqfield":  "%sと確認用パスワードが違います。",
		"json.invalid":      "jsonが不正です",
		"login.failed":      "ログインに失敗しました",
	},
}
