package taskRoom

import (
	"context"
	"encoding/json"

	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/NoobforAl/real_time_chat_application/src/entity"
	taskTypes "github.com/NoobforAl/real_time_chat_application/src/tasks/tasks_type"
	"github.com/hibiken/asynq"
)

type RoomSaveProcess struct {
	store contract.StoreRoomsAndCache
	log   contract.Logger
}

func NewRoomSaveTask(room entity.Room) (*asynq.Task, error) {
	msg, err := json.Marshal(room)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(taskTypes.TypeRoomSave, msg), nil
}

func NewRoomSaveProcess(store contract.StoreRoomsAndCache, log contract.Logger) *RoomSaveProcess {
	return &RoomSaveProcess{store: store, log: log}
}

func (msp *RoomSaveProcess) ProcessTask(ctx context.Context, t *asynq.Task) error {
	dataMessage := t.Payload()
	msp.log.Debugf("process save message, message content: %v", dataMessage)

	var Room entity.Room
	err := json.Unmarshal(dataMessage, &Room)
	if err != nil {
		return err
	}

	// set in queue for show with websocket to online user
	err = msp.store.SetRooms(ctx, []*entity.Room{&Room})
	if err != nil {
		return err
	}

	_, err = msp.store.CreateRoom(ctx, Room)
	return err
}
