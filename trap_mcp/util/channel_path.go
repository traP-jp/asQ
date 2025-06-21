package util

import (
	"fmt"

	"github.com/traPtitech/go-traq"
)

func ConvertChannelNameToPath(original *traq.ChannelList) (traq.ChannelList, error) {
	if original == nil {
		return traq.ChannelList{}, fmt.Errorf("original channel list is nil")
	}
	idToName := make(map[string]string)
	idToParentId := make(map[string]string)
	for _, channel := range original.Public {
		idToName[channel.Id] = channel.Name
		if parentId, ok := channel.GetParentIdOk(); ok && parentId != nil {
			idToParentId[channel.Id] = *parentId
		}
	}
	revised := make([]traq.Channel, 0, len(original.Public))
	for _, channel := range original.Public {
		path, err := getPathInner(&idToName, &idToParentId, channel.Id)
		if err != nil {
			return traq.ChannelList{}, fmt.Errorf("failed to convert channel name to path: %w", err)
		}
		revisedChannel := channel
		revisedChannel.Name = path
		revised = append(revised, revisedChannel)
	}
	return traq.ChannelList{
		Public: revised,
		Dm:     original.Dm,
	}, nil
}

func getPathInner(idToName map[string]string, idToParentId map[string]string, channelId string) (string, error) {
	if originalName, ok := idToName[channelId]; ok {
		if parentId, ok := idToParentId[channelId]; ok && parentId != "" {
			parentPath, err := getPathInner(idToName, idToParentId, parentId)
			if err != nil {
				return "", err
			}
			return parentPath + "/" + originalName, nil
		}
		return originalName, nil
	}
	return "", fmt.Errorf("channel ID %s not found in idToName map", channelId)

}
