package compute

import (
	"github.com/rancher/agent/model"
	"github.com/rancher/agent/utilities/utils"
	"github.com/rancher/agent/utilities/constants"
	"github.com/docker/engine-api/client"
	"github.com/pkg/errors"
)

func IsInstanceActive(instance model.Instance, host model.Host, client *client.Client) (bool, error) {
	if utils.IsNoOp(instance.ProcessData) {
		return true, nil
	}

	container, err := utils.GetContainer(client, instance, false)
	if err != nil {
		if utils.IsContainerNotFoundError(err) {
			return false, nil
		}
		return false, errors.Wrap(err, constants.IsInstanceActiveError)
	}
	return isRunning(client, container)
}

func IsInstanceInactive(instance model.Instance, client *client.Client) (bool, error) {
	if utils.IsNoOp(instance.ProcessData) {
		return true, nil
	}

	container, err := utils.GetContainer(client, instance, false)
	if err != nil {
		if !utils.IsContainerNotFoundError(err) {
			return false, errors.Wrap(err, constants.IsInstanceInactiveError)
		}
	}
	return isStopped(client, container)
}

func IsInstanceRemoved(instance model.Instance, dockerClient *client.Client) (bool, error) {
	con, err := utils.GetContainer(dockerClient, instance, false)
	if err != nil {
		if utils.IsContainerNotFoundError(err) {
			return true, nil
		}
		return false, errors.Wrap(err, constants.IsInstanceRemovedError)
	}
	return con.ID == "", nil
}