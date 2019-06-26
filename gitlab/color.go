package gitlab

// ColorGitlab represents GitLab color
const ColorGitlab = "#fc6d26" // https://brandcolors.net/b/gitlab

func ciStatusColor(status string) string {
	// https://gitlab.com/gitlab-org/gitlab-ce/blob/v11.11.4/app/assets/stylesheets/pages/status.scss
	switch status {
	case "failed":
		return "#db3b21"
	case "success":
		return "#1aaa55"
	case "canceled", "disabled", "scheduled", "manual":
		return "#2e2e2e"
	case "pending", "failed-with-warnings", "success-with-warnings":
		return "#fc9403"
	case "info", "preparing", "running":
		return "#1f78d1"
	case "created", "skipped":
		return "#707070"
	default:
		return ColorGitlab
	}
}
