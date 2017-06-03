package config

import (
	"fmt"
	"os"
)

func GetPasswordHashForUser(user string) string {
	// TODO impl. user dependent hash retrieval
	return os.Getenv("GOMIC_BASIC_AUTH_PASS_HASH")
}

func GetTwitterRepeatConsumerKey() string {
	s := os.Getenv("TWITTER_REPEAT_CONSUMER_KEY")
	fmt.Println(s)
	return s
}

func GetTwitterRepeatConsumerSecret() string {
	s := os.Getenv("TWITTER_REPEAT_CONSUMER_SECRET")
	fmt.Println(s)
	return s
}

func GetTwitterRepeatAccessToken() string {
	s := os.Getenv("TWITTER_REPEAT_ACCESS_TOKEN")
	fmt.Println(s)
	return s
}

func GetTwitterRepeatAccessTokenSecret() string {
	s := os.Getenv("TWITTER_REPEAT_ACCESS_TOKEN_SECRET")
	fmt.Println(s)
	return s
}

func GetTwitterConsumerKey() string {
	s := os.Getenv("TWITTER_ORIGIN_CONSUMER_KEY")
	fmt.Println(s)
	return s
}

func GetTwitterConsumerSecret() string {
	s := os.Getenv("TWITTER_ORIGIN_CONSUMER_SECRET")
	fmt.Println(s)
	return s
}

func GetTwitterAccessToken() string {
	s := os.Getenv("TWITTER_ORIGIN_ACCESS_TOKEN")
	fmt.Println(s)
	return s
}

func GetTwitterAccessTokenSecret() string {
	s := os.Getenv("TWITTER_ORIGIN_ACCESS_TOKEN_SECRET")
	fmt.Println(s)
	return s
}

func GetTumblrCallbackUrl() string {
	return os.Getenv("GOMIC_TUMBLR_CALLBACK_URL")
}

func GetTumblrBlogName() string {
	return os.Getenv("GOMIC_TUMBLR_BLOG_NAME")
}

func GetTumblrConsumerKey() string {
	return os.Getenv("GOMIC_TUMBLR_CONSUMER_KEY")
}

func GetTumblrConsumerSecret() string {
	return os.Getenv("GOMIC_TUMBLR_CONSUMER_SECRET")
}

func GetTumblrToken() string {
	return os.Getenv("GOMIC_TUMBLR_TOKEN")
}

func GetTumblrTokenSecret() string {
	return os.Getenv("GOMIC_TUMBLR_TOKEN_SECRET")
}

func GetFacebookPageId() string {
	return os.Getenv("FB_DEVABODE_ID")
}

func GetFacebookApplicationId() string {
	return os.Getenv("FB_APPLICATION_ID")
}

func GetFacebookApplicationSecret() string {
	return os.Getenv("FB_APPLICATION_SECRET")
}

func GetFacebookCallbackUrl() string {
	return os.Getenv("FB_CALLBACK_URL")
}

func IsProd() bool {
	return os.Getenv("GOMIC_STAGE") == "prod"
}

func GetDsn() string {
	user := os.Getenv("DB_GOMIC_USER")
	pass := os.Getenv("DB_GOMIC_PASS")
	name := os.Getenv("DB_GOMIC_NAME")
	host := os.Getenv("DB_GOMIC_HOST")
	return fmt.Sprintf("%s:%s@%s/%s", user, pass, host, name)
}

func GetTlsPaths() (string, string) {
	cert := os.Getenv("TLS_CERT_PATH")
	key := os.Getenv("TLS_KEY_PATH")
	return cert, key
}
