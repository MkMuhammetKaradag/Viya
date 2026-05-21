package domain

import (
	"context"

	"github.com/google/uuid"
)

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *Trip) (uuid.UUID, error)
	AddWaypoint(ctx context.Context, wp *Waypoint) (uuid.UUID, error)
	AddWaypointPhotos(ctx context.Context, waypointID uuid.UUID, photoURLs []string) error
	AddWaypointPhotoWithTags(ctx context.Context, wpID uuid.UUID, photoURL string, tags []Tag) error
	SearchCategories(ctx context.Context, query string) ([]Category, error)
	// GetTripByID(ctx context.Context, tripID uuid.UUID) (*Trip, error)

	GetTripWithWaypointsAndPhotos(ctx context.Context, tripID, currentUserID uuid.UUID) (*Trip, error)
	IncrementUniqueView(ctx context.Context, tripID, userID uuid.UUID) error

	GetMeTrips(ctx context.Context, userID uuid.UUID, page, limit int) ([]TripSummary, error)
	GetUserTrips(ctx context.Context, currentUserID, targetUserID uuid.UUID, page, limit int) ([]TripSummary, error)
	GetLikedTrips(ctx context.Context, userID uuid.UUID, page, limit int) ([]TripSummary, error)
	GetExploreTrips(ctx context.Context, userID uuid.UUID, limit, offset int) ([]TripExploreDTO, error)
	GetHomeFeedTrips(ctx context.Context, userID uuid.UUID, limit, offset int) ([]TripExploreDTO, error)
	GetTripStatus(ctx context.Context, tripID uuid.UUID) (*TripStatusDTO, error)
	ToggleTripLike(ctx context.Context, tripID, userID uuid.UUID) (bool, error)

	DeleteWaypoint(ctx context.Context, waypointID uuid.UUID) error
	ReorderWaypoints(ctx context.Context, wpID uuid.UUID, index int) error
	GetWaypointByID(ctx context.Context, id uuid.UUID) (*Waypoint, error)
	UpdateWaypoint(ctx context.Context, wp *Waypoint) error

	CreateUser(ctx context.Context, id uuid.UUID, username, email string) error
	UpdateUserSocialInfo(ctx context.Context, id uuid.UUID, isPrivate *bool, avatarURL *string) error

	GetTripByIDForAI(ctx context.Context, id uuid.UUID) (*Trip, error)
	UpdateTripEmbedding(ctx context.Context, id uuid.UUID, vector []float32) error
	UpdateUserInterest(ctx context.Context, userID uuid.UUID, tripID uuid.UUID, weight float32) error

	UpsertLocalFollow(ctx context.Context, followerID, followingID uuid.UUID, status string) error
	CheckFollowStatus(ctx context.Context, followerID, followedID uuid.UUID) (bool, error)

	UpsertLocalBlock(ctx context.Context, blockerID, blockedID uuid.UUID) error

	CreateComment(ctx context.Context, comment *Comment) (uuid.UUID, error)
	GetTripComments(ctx context.Context, viewerID uuid.UUID, tripID uuid.UUID, page, limit int) ([]Comment, error)
	GetCommentReplies(ctx context.Context, parentID uuid.UUID, page, limit int) ([]Comment, error)

	ForkTrip(ctx context.Context, originalTripID uuid.UUID, forkUserID uuid.UUID) (uuid.UUID, error)
	Close() error
}
