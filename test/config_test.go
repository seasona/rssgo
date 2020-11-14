package test

import(
	"testing"

	"github.com/seasona/rssgo/internal"
)

func TestLoadConfig(t *testing.T){
	url := "../config.json";
	
	config:= internal.Config{}

	config.LoadConfig(url)
}