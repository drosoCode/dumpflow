package importer

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	xmlparser "github.com/tamerh/xml-stream-parser"

	"github.com/drosocode/dumpflow/internal/database"
)

func skipErr[T any](data T, err error) T {
	return data
}

var CopySize int

func init() {
	CopySize = 500
}

func importFile(db *database.Queries, file io.ReadCloser, name string, size int64, st *ImportStatus, wg *sync.WaitGroup) {
	fmt.Printf("starting import of %s\n", name)

	switch name {
	case "Badges.xml":
		importBadges(db, file, size, st)
	case "Comments.xml":
		importComments(db, file, size, st)
	case "PostHistory.xml":
		importPostHistory(db, file, size, st)
	case "PostLinks.xml":
		importPostLinks(db, file, size, st)
	case "Posts.xml":
		importPosts(db, file, size, st)
	case "Tags.xml":
		importTags(db, file, size, st)
	case "Users.xml":
		importUsers(db, file, size, st)
	case "Votes.xml":
		importVotes(db, file, size, st)
	}

	wg.Done()
	fmt.Printf("finished import of %s\n", name)
}

func importBadges(db *database.Queries, file io.ReadCloser, size int64, status *ImportStatus) {
	status.Badges.Total = size

	ctx := context.Background()
	br := bufio.NewReaderSize(file, 1024*1024)
	parser := xmlparser.NewXMLParser(br, "row")
	data := []database.AddBadgeParams{}
	i := 0

	for xml := range parser.Stream() {
		data = append(data, database.AddBadgeParams{
			ID:       skipErr(strconv.ParseInt(xml.Attrs["Id"], 10, 64)),
			UserID:   skipErr(strconv.ParseInt(xml.Attrs["UserId"], 10, 64)),
			Name:     xml.Attrs["Name"],
			Date:     skipErr(time.Parse("2006-01-02T15:04:05.000", xml.Attrs["Date"])),
			Class:    int32(skipErr(strconv.ParseInt(xml.Attrs["Class"], 10, 32))),
			TagBased: xml.Attrs["TagBased"] == "true",
		})
		i++
		status.Badges.Current = int64(parser.TotalReadSize)
		if i == CopySize {
			db.AddBadge(ctx, data)
			i = 0
			data = []database.AddBadgeParams{}
		}
	}
	db.AddBadge(ctx, data)
}

func importComments(db *database.Queries, file io.ReadCloser, size int64, status *ImportStatus) {
	status.Comments.Total = size

	ctx := context.Background()
	br := bufio.NewReaderSize(file, 1024*1024)
	parser := xmlparser.NewXMLParser(br, "row")
	data := []database.AddCommentParams{}
	i := 0

	for xml := range parser.Stream() {
		data = append(data, database.AddCommentParams{
			ID:             skipErr(strconv.ParseInt(xml.Attrs["Id"], 10, 64)),
			PostID:         skipErr(strconv.ParseInt(xml.Attrs["PostId"], 10, 64)),
			Score:          int32(skipErr(strconv.ParseInt(xml.Attrs["Score"], 10, 32))),
			Text:           xml.Attrs["Text"],
			CreationDate:   skipErr(time.Parse("2006-01-02T15:04:05.000", xml.Attrs["CreationDate"])),
			UserID:         skipErr(strconv.ParseInt(xml.Attrs["UserId"], 10, 64)),
			ContentLicense: xml.Attrs["ContentLicense"],
		})
		i++
		status.Comments.Current = int64(parser.TotalReadSize)
		if i == CopySize {
			db.AddComment(ctx, data)
			i = 0
			data = []database.AddCommentParams{}
		}
	}
	db.AddComment(ctx, data)
}

func importPostHistory(db *database.Queries, file io.ReadCloser, size int64, status *ImportStatus) {
	status.PostHistory.Total = size

	ctx := context.Background()
	br := bufio.NewReaderSize(file, 1024*1024)
	parser := xmlparser.NewXMLParser(br, "row")
	data := []database.AddPostHistoryParams{}
	i := 0

	for xml := range parser.Stream() {
		data = append(data, database.AddPostHistoryParams{
			ID:                skipErr(strconv.ParseInt(xml.Attrs["Id"], 10, 64)),
			PostHistoryTypeID: int32(skipErr(strconv.ParseInt(xml.Attrs["PostHistoryTypeId"], 10, 32))),
			PostID:            skipErr(strconv.ParseInt(xml.Attrs["PostId"], 10, 64)),
			RevisionGuid:      xml.Attrs["RevisionGUID"],
			CreationDate:      skipErr(time.Parse("2006-01-02T15:04:05.000", xml.Attrs["CreationDate"])),
			UserID:            skipErr(strconv.ParseInt(xml.Attrs["UserId"], 10, 64)),
			Comment:           xml.Attrs["Comment"],
			Text:              xml.Attrs["Text"],
			ContentLicense:    xml.Attrs["ContentLicense"],
		})
		i++
		status.PostHistory.Current = int64(parser.TotalReadSize)
		if i == CopySize {
			db.AddPostHistory(ctx, data)
			i = 0
			data = []database.AddPostHistoryParams{}
		}
	}
	db.AddPostHistory(ctx, data)
}

func importPostLinks(db *database.Queries, file io.ReadCloser, size int64, status *ImportStatus) {
	status.PostLinks.Total = size

	ctx := context.Background()
	br := bufio.NewReaderSize(file, 1024*1024)
	parser := xmlparser.NewXMLParser(br, "row")
	data := []database.AddPostLinkParams{}
	i := 0

	for xml := range parser.Stream() {
		data = append(data, database.AddPostLinkParams{
			ID:            skipErr(strconv.ParseInt(xml.Attrs["Id"], 10, 64)),
			CreationDate:  skipErr(time.Parse("2006-01-02T15:04:05.000", xml.Attrs["CreationDate"])),
			PostID:        skipErr(strconv.ParseInt(xml.Attrs["PostId"], 10, 64)),
			RelatedPostID: skipErr(strconv.ParseInt(xml.Attrs["RelatedPostId"], 10, 64)),
			LinkTypeID:    int32(skipErr(strconv.ParseInt(xml.Attrs["LinkTypeId"], 10, 32))),
		})
		i++
		status.PostLinks.Current = int64(parser.TotalReadSize)
		if i == CopySize {
			db.AddPostLink(ctx, data)
			i = 0
			data = []database.AddPostLinkParams{}
		}
	}
	db.AddPostLink(ctx, data)
}

func importPosts(db *database.Queries, file io.ReadCloser, size int64, status *ImportStatus) {
	status.Posts.Total = size

	ctx := context.Background()
	br := bufio.NewReaderSize(file, 1024*1024)
	parser := xmlparser.NewXMLParser(br, "row")
	data := []database.AddPostParams{}
	i := 0

	for xml := range parser.Stream() {
		data = append(data, database.AddPostParams{
			ID:               skipErr(strconv.ParseInt(xml.Attrs["Id"], 10, 64)),
			PostTypeID:       int32(skipErr(strconv.ParseInt(xml.Attrs["PostTypeId"], 10, 32))),
			ParentID:         skipErr(strconv.ParseInt(xml.Attrs["ParentId"], 10, 64)),
			AcceptedAnswerID: int32(skipErr(strconv.ParseInt(xml.Attrs["AcceptedAnswerId"], 10, 32))),
			CreationDate:     skipErr(time.Parse("2006-01-02T15:04:05.000", xml.Attrs["CreationDate"])),
			ClosedDate:       skipErr(time.Parse("2006-01-02T15:04:05.000", xml.Attrs["ClosedDate"])),
			Score:            int32(skipErr(strconv.ParseInt(xml.Attrs["Score"], 10, 32))),
			ViewCount:        int32(skipErr(strconv.ParseInt(xml.Attrs["ViewCount"], 10, 32))),
			Body:             xml.Attrs["Body"],
			Tags:             xml.Attrs["Tags"],
			AnswerCount:      int32(skipErr(strconv.ParseInt(xml.Attrs["AnswerCount"], 10, 32))),
			CommentCount:     int32(skipErr(strconv.ParseInt(xml.Attrs["CommentCount"], 10, 32))),
			FavoriteCount:    int32(skipErr(strconv.ParseInt(xml.Attrs["FavoriteCount"], 10, 32))),
			ContentLicense:   xml.Attrs["ContentLicense"],
		})
		i++
		status.Posts.Current = int64(parser.TotalReadSize)
		if i == CopySize {
			db.AddPost(ctx, data)
			i = 0
			data = []database.AddPostParams{}
		}
	}
	db.AddPost(ctx, data)
}

func importTags(db *database.Queries, file io.ReadCloser, size int64, status *ImportStatus) {
	status.Tags.Total = size

	ctx := context.Background()
	br := bufio.NewReaderSize(file, 1024*1024)
	parser := xmlparser.NewXMLParser(br, "row")
	data := []database.AddTagParams{}
	i := 0

	for xml := range parser.Stream() {
		data = append(data, database.AddTagParams{
			ID:              skipErr(strconv.ParseInt(xml.Attrs["Id"], 10, 64)),
			TagName:         xml.Attrs["TagName"],
			Count:           int32(skipErr(strconv.ParseInt(xml.Attrs["Count"], 10, 32))),
			IsRequired:      xml.Attrs["IsRequired"] == "true",
			IsModeratorOnly: xml.Attrs["IsModeratorOnly"] == "true",
			WikiPostID:      int32(skipErr(strconv.ParseInt(xml.Attrs["WikiPostId"], 10, 32))),
			ExcerptPostID:   int32(skipErr(strconv.ParseInt(xml.Attrs["ExcerptPostId"], 10, 32))),
		})
		i++
		status.Tags.Current = int64(parser.TotalReadSize)
		if i == CopySize {
			db.AddTag(ctx, data)
			i = 0
			data = []database.AddTagParams{}
		}
	}
	db.AddTag(ctx, data)
}

func importUsers(db *database.Queries, file io.ReadCloser, size int64, status *ImportStatus) {
	status.Users.Total = size

	ctx := context.Background()
	br := bufio.NewReaderSize(file, 1024*1024)
	parser := xmlparser.NewXMLParser(br, "row")
	data := []database.AddUserParams{}
	i := 0

	for xml := range parser.Stream() {
		data = append(data, database.AddUserParams{
			ID:              skipErr(strconv.ParseInt(xml.Attrs["Id"], 10, 64)),
			Reputation:      int32(skipErr(strconv.ParseInt(xml.Attrs["Reputation"], 10, 32))),
			CreationDate:    skipErr(time.Parse("2006-01-02T15:04:05.000", xml.Attrs["CreationDate"])),
			DisplayName:     xml.Attrs["DisplayName"],
			LastAccessDate:  skipErr(time.Parse("2006-01-02T15:04:05.000", xml.Attrs["LastAccessDate"])),
			Location:        xml.Attrs["Location"],
			WebsiteUrl:      xml.Attrs["WebsiteUrl"],
			AboutMe:         xml.Attrs["AboutMe"],
			Views:           int32(skipErr(strconv.ParseInt(xml.Attrs["Views"], 10, 32))),
			Upvotes:         int32(skipErr(strconv.ParseInt(xml.Attrs["UpVotes"], 10, 32))),
			Downvotes:       int32(skipErr(strconv.ParseInt(xml.Attrs["DownVotes"], 10, 32))),
			AccountID:       skipErr(strconv.ParseInt(xml.Attrs["AccountId"], 10, 64)),
			ProfileImageUrl: xml.Attrs["ProfileImageUrl"],
		})
		i++
		status.Users.Current = int64(parser.TotalReadSize)
		if i == CopySize {
			db.AddUser(ctx, data)
			i = 0
			data = []database.AddUserParams{}
		}
	}
	db.AddUser(ctx, data)
}

func importVotes(db *database.Queries, file io.ReadCloser, size int64, status *ImportStatus) {
	status.Votes.Total = size
	ctx := context.Background()
	br := bufio.NewReaderSize(file, 1024*1024)
	parser := xmlparser.NewXMLParser(br, "row")
	data := []database.AddVoteParams{}
	i := 0

	for xml := range parser.Stream() {
		data = append(data, database.AddVoteParams{
			ID:           skipErr(strconv.ParseInt(xml.Attrs["Id"], 10, 64)),
			PostID:       skipErr(strconv.ParseInt(xml.Attrs["PostId"], 10, 64)),
			VoteTypeID:   int32(skipErr(strconv.ParseInt(xml.Attrs["VoteTypeId"], 10, 32))),
			CreationDate: skipErr(time.Parse("2006-01-02T15:04:05.000", xml.Attrs["CreationDate"])),
		})
		i++
		status.Votes.Current = int64(parser.TotalReadSize)
		if i == CopySize {
			db.AddVote(ctx, data)
			i = 0
			data = []database.AddVoteParams{}
		}
	}
	db.AddVote(ctx, data)
}
