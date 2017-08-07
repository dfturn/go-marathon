package marathon

type Service interface {
	Versions() (*ApplicationVersions, error)

	Update(force bool) error

	Delete(force bool) (*DeploymentID, error)

	Scale(instances int, force bool) (*DeploymentID, error)

	Restart(force bool) (*DeploymentID, error)

	ExistsAndRunning() bool
}
