package main

import "github.com/circulohealth/sonar-backend/packages/router/pkg/forward"

func Handler(dto forward.ForwarderBroadcastDTO) error {
	forwarder, err := forward.NewForwarderFromBroadcast(dto)

	if err != nil {
		return err
	}

	err = forwarder.Forward()

	if err != nil {
		return err
	}

	return nil
}
