package centrifugo

import (
	"context"
	"encoding/json"
	"fmt"
	"github/socialforge/internal/infra/contextpool"
	"time"

	"github.com/centrifugal/gocent/v3"
	"go.uber.org/zap"
)

/**
 * IsUp checks if the Centrifugo client is up and running
 * @return {bool} - True if the client is up, false otherwise
 */
func (c *CentrifugoClient) IsUp() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isUp
}

/**
 * PublishMessage publishes a message to a specified Centrifugo channel
 * @param {context.Context} ctx - Context for request cancellation
 * @param {string} channel - The channel to publish the message to
 * @param {interface{}} data - The data to publish (will be JSON marshaled)
 * @return {error} - Error if the publish operation fails
 */
func (c *CentrifugoClient) PublishMessage(ctx context.Context, channel string, data interface{}) error {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 5*time.Second)
	defer cancel()

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	_, err = c.client.Publish(ctx, channel, dataBytes)
	if err != nil {
		c.logger.Error("Failed to publish message to Centrifugo",
			zap.String("channel", channel),
			zap.Error(err),
		)
		return fmt.Errorf("failed to publish to channel %s: %w", channel, err)
	}

	c.logger.Debug("Message published to Centrifugo",
		zap.String("channel", channel),
	)
	return nil
}

/**
 * PublishToUser publishes a message to a user-specific Centrifugo channel
 * @param {context.Context} ctx - Context for request cancellation
 * @param {string} userID - The ID of the user to publish to
 * @param {interface{}} data - The data to publish (will be JSON marshaled)
 * @return {error} - Error if the publish operation fails
 */
func (c *CentrifugoClient) PublishToUser(ctx context.Context, userID string, data interface{}) error {
	channel := fmt.Sprintf("user:%s", userID)
	return c.PublishMessage(ctx, channel, data)
}

/**
 * PublishToConversation publishes a message to a conversation-specific Centrifugo channel
 * @param {context.Context} ctx - Context for request cancellation
 * @param {string} conversationID - The ID of the conversation to publish to
 * @param {interface{}} data - The data to publish (will be JSON marshaled)
 * @return {error} - Error if the publish operation fails
 */
func (c *CentrifugoClient) PublishToConversation(ctx context.Context, conversationID string, data interface{}) error {
	channel := fmt.Sprintf("conversation:%s", conversationID)
	return c.PublishMessage(ctx, channel, data)
}

/**
 * PublishToTenant publishes a message to a tenant-specific Centrifugo channel
 * @param {context.Context} ctx - Context for request cancellation
 * @param {string} tenantID - The ID of the tenant to publish to
 * @param {interface{}} data - The data to publish (will be JSON marshaled)
 * @return {error} - Error if the publish operation fails
 */
func (c *CentrifugoClient) PublishToTenant(ctx context.Context, tenantID string, data interface{}) error {
	channel := fmt.Sprintf("tenant:%s", tenantID)
	return c.PublishMessage(ctx, channel, data)
}

/**
 * PublishToDivision publishes a message to a division-specific Centrifugo channel
 * @param {context.Context} ctx - Context for request cancellation
 * @param {string} divisionID - The ID of the division to publish to
 * @param {interface{}} data - The data to publish (will be JSON marshaled)
 * @return {error} - Error if the publish operation fails
 */
func (c *CentrifugoClient) PublishToDivision(ctx context.Context, divisionID string, data interface{}) error {
	channel := fmt.Sprintf("division:%s", divisionID)
	return c.PublishMessage(ctx, channel, data)
}

/**
 * BroadcastNewMessage broadcasts a new message event to a conversation-specific Centrifugo channel
 * @param {context.Context} ctx - Context for request cancellation
 * @param {string} conversationID - The ID of the conversation to broadcast to
 * @param {map[string]interface{}} message - The message payload to broadcast
 * @return {error} - Error if the publish operation fails
 */
func (c *CentrifugoClient) BroadcastNewMessage(ctx context.Context, conversationID string, message map[string]interface{}) error {
	payload := map[string]interface{}{
		"type":    "new_message",
		"message": message,
	}
	return c.PublishToConversation(ctx, conversationID, payload)
}

/**
 * BroadcastTypingIndicator broadcasts a typing indicator event to a conversation-specific Centrifugo channel
 * @param {context.Context} ctx - Context for request cancellation
 * @param {string} conversationID - The ID of the conversation to broadcast to
 * @param {string} userID - The ID of the user typing
 * @param {bool} isTyping - Whether the user is typing or not
 * @return {error} - Error if the publish operation fails
 */
func (c *CentrifugoClient) BroadcastTypingIndicator(ctx context.Context, conversationID, userID string, isTyping bool) error {
	payload := map[string]interface{}{
		"type":      "typing",
		"user_id":   userID,
		"is_typing": isTyping,
	}
	return c.PublishToConversation(ctx, conversationID, payload)
}

/**
 * BroadcastMessageRead broadcasts a message read event to a conversation-specific Centrifugo channel
 * @param {context.Context} ctx - Context for request cancellation
 * @param {string} conversationID - The ID of the conversation to broadcast to
 * @param {string} userID - The ID of the user who read the messages
 * @param {[]string} messageIDs - The IDs of the messages that were read
 * @return {error} - Error if the publish operation fails
 */
func (c *CentrifugoClient) BroadcastMessageRead(ctx context.Context, conversationID, userID string, messageIDs []string) error {
	payload := map[string]interface{}{
		"type":        "message_read",
		"user_id":     userID,
		"message_ids": messageIDs,
	}
	return c.PublishToConversation(ctx, conversationID, payload)
}

/**
 * BroadcastConversationUpdate broadcasts a conversation update event to a conversation-specific Centrifugo channel
 * @param {context.Context} ctx - Context for request cancellation
 * @param {string} conversationID - The ID of the conversation to broadcast to
 * @param {map[string]interface{}} updates - The conversation updates payload to broadcast
 * @return {error} - Error if the publish operation fails
 */
func (c *CentrifugoClient) BroadcastConversationUpdate(ctx context.Context, conversationID string, updates map[string]interface{}) error {
	payload := map[string]interface{}{
		"type":    "conversation_update",
		"updates": updates,
	}
	return c.PublishToConversation(ctx, conversationID, payload)
}

/**
 * BroadcastAgentAssigned broadcasts an agent assigned event to a conversation-specific Centrifugo channel
 * @param {context.Context} ctx - Context for request cancellation
 * @param {string} conversationID - The ID of the conversation to broadcast to
 * @param {string} agentID - The ID of the agent assigned to the conversation
 * @return {error} - Error if the publish operation fails
 */
func (c *CentrifugoClient) BroadcastAgentAssigned(ctx context.Context, conversationID, agentID string) error {
	payload := map[string]interface{}{
		"type":     "agent_assigned",
		"agent_id": agentID,
	}
	return c.PublishToConversation(ctx, conversationID, payload)
}

/**
 * BroadcastConversationClosed broadcasts a conversation closed event to a conversation-specific Centrifugo channel
 * @param {context.Context} ctx - Context for request cancellation
 * @param {string} conversationID - The ID of the conversation to broadcast to
 * @param {string} closedBy - The ID of the user or system that closed the conversation
 * @return {error} - Error if the publish operation fails
 */
func (c *CentrifugoClient) BroadcastConversationClosed(ctx context.Context, conversationID, closedBy string) error {
	payload := map[string]interface{}{
		"type":      "conversation_closed",
		"closed_by": closedBy,
		"closed_at": time.Now().Unix(),
	}
	return c.PublishToConversation(ctx, conversationID, payload)
}

/**
 * BroadcastConversationReopened broadcasts a conversation reopened event to a conversation-specific Centrifugo channel
 * @param {context.Context} ctx - Context for request cancellation
 * @param {string} conversationID - The ID of the conversation to broadcast to
 * @param {string} reopenedBy - The ID of the user or system that reopened the conversation
 * @return {error} - Error if the publish operation fails
 */
func (c *CentrifugoClient) BroadcastConversationReopened(ctx context.Context, conversationID, reopenedBy string) error {
	payload := map[string]interface{}{
		"type":        "conversation_reopened",
		"reopened_by": reopenedBy,
		"reopened_at": time.Now().Unix(),
	}
	return c.PublishToConversation(ctx, conversationID, payload)
}

/**
 * Presence retrieves the presence information for a given Centrifugo channel
 * @param {context.Context} ctx - Context for request cancellation
 * @param {string} channel - The Centrifugo channel to query presence for
 * @return {*gocent.PresenceResult} - The presence result containing client IDs and their data
 * @return {error} - Error if the presence operation fails
 */
func (c *CentrifugoClient) Presence(ctx context.Context, channel string) (*gocent.PresenceResult, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 5*time.Second)
	defer cancel()

	result, err := c.client.Presence(ctx, channel)
	if err != nil {
		return nil, fmt.Errorf("failed to get presence for channel %s: %w", channel, err)
	}
	return &result, nil
}

/**
 * History retrieves the message history for a given Centrifugo channel
 * @param {context.Context} ctx - Context for request cancellation
 * @param {string} channel - The Centrifugo channel to query history for
 * @param {int} limit - The maximum number of messages to retrieve
 * @return {*gocent.HistoryResult} - The history result containing messages
 * @return {error} - Error if the history operation fails
 */
func (c *CentrifugoClient) History(ctx context.Context, channel string, limit int) (*gocent.HistoryResult, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 5*time.Second)
	defer cancel()

	result, err := c.client.History(ctx, channel, gocent.WithLimit(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to get history for channel %s: %w", channel, err)
	}
	return &result, nil
}

/**
 * Channels retrieves the list of active Centrifugo channels
 * @param {context.Context} ctx - Context for request cancellation
 * @return {[]string} - The list of active channel names
 * @return {error} - Error if the channels operation fails
 */
func (c *CentrifugoClient) Channels(ctx context.Context) ([]string, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 5*time.Second)
	defer cancel()

	result, err := c.client.Channels(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get channels: %w", err)
	}
	channels := make([]string, 0, len(result.Channels))
	for ch := range result.Channels {
		channels = append(channels, ch)
	}
	return channels, nil
}

/**
 * Disconnect disconnects a user from all Centrifugo channels
 * @param {context.Context} ctx - Context for request cancellation
 * @param {string} userID - The ID of the user to disconnect
 * @return {error} - Error if the disconnect operation fails
 */
func (c *CentrifugoClient) Disconnect(ctx context.Context, userID string, opts ...gocent.DisconnectOption) error {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 5*time.Second)
	defer cancel()

	err := c.client.Disconnect(ctx, userID, opts...)
	if err != nil {
		return fmt.Errorf("failed to disconnect user %s: %w", userID, err)
	}
	return nil
}

/**
 * Refresh refreshes the token for a user in Centrifugo
 * @param {context.Context} ctx - Context for request cancellation
 * @param {string} userID - The ID of the user to refresh
 * @return {error} - Error if the refresh operation fails
 */
func (c *CentrifugoClient) Refresh(ctx context.Context, userID string, opts ...gocent.DisconnectOption) error {
	_, cancel := contextpool.WithTimeoutIfNone(ctx, 5*time.Second)
	defer cancel()

	err := c.client.Pipe().AddDisconnect(userID, opts...)
	if err != nil {
		return fmt.Errorf("failed to refresh user %s: %w", userID, err)
	}
	return nil
}

/**
 * Close closes the Centrifugo client connection
 * @return {error} - Error if the close operation fails
 */
func (c *CentrifugoClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.isUp = false
	c.logger.Info("Centrifugo client closed")
	return nil
}
