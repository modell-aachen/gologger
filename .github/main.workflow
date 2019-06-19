workflow "Build and Publish" {
  on = "push"
  resolves = ["Publish"]
}

action "Build" {
  uses = "cedrickring/golang-action@1.3.0"
  runs = "sh"
  args = ["build.sh"]
}

action "Publish" {
  uses = "elgohr/Github-Hub-Action@1.0"
  args = "release create $(date +%Y%m%d%H%M%S) -a cmd/logstore/logstore -a cmd/logtail/logtail -a cmd/logreport/logreport -a cmd/logdump/logdump"
  secrets = ["GITHUB_TOKEN"]
  needs = ["Build"]
}