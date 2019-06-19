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
  needs = ["Build"]
  secrets = ["GITHUB_TOKEN"]
}
