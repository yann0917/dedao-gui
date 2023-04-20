package services

type NotesUser struct {
	Id          string `json:"id"`
	Uid         int    `json:"uid"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	Follow      int    `json:"follow"`
	IsV         int    `json:"isV"`
	Slogan      string `json:"slogan"`
	Vinfo       string `json:"Vinfo"`
	StudentId   int    `json:"student_id"`
	IsPoster    bool   `json:"is_poster"`
	Qrcode      string `json:"qrcode"`
	LogId       string `json:"log_id"`
	LogType     string `json:"log_type"`
	UidHazy     string `json:"uid_hazy"`
	NoteIdHazy  string `json:"note_id_hazy"`
	Role        string `json:"role"`
	Attribution string `json:"attribution"`
}

type Comment struct {
	AttachmentType          int            `json:"attachment_type"`
	AuditState              int            `json:"audit_state"`
	Class                   int            `json:"class"`
	CommentReply            string         `json:"comment_reply"`
	CommentReplyTime        int            `json:"comment_reply_time"`
	Content                 string         `json:"content"`
	ContentType             int            `json:"content_type"`
	CreateTime              int            `json:"create_time"`
	CurrentState            int            `json:"current_state"`
	Ddurl                   Ddurl          `json:"ddurl"`
	DetailId                int            `json:"detail_id"`
	Extra                   NotesExtra     `json:"extra"`
	FeedId                  int64          `json:"feed_id"`
	Highlights              []interface{}  `json:"highlights"`
	IsFromMe                int            `json:"is_from_me"`
	CanEdit                 bool           `json:"can_edit"`
	IsLike                  bool           `json:"is_like"`
	IsOpenLedgers           bool           `json:"is_open_ledgers"`
	IsPermission            bool           `json:"is_permission"`
	IsReposted              bool           `json:"is_reposted"`
	IsUpmost                bool           `json:"is_upmost"`
	LevelPermission         bool           `json:"level_permission"`
	Lesson                  Lesson         `json:"lesson"`
	Level                   int            `json:"level"`
	LevelType               int            `json:"level_type"`
	LogId                   string         `json:"log_id"`
	LogType                 string         `json:"log_type"`
	Note                    string         `json:"note"`
	NoteId                  int64          `json:"note_id"`
	NoteIdHazy              string         `json:"note_id_hazy"`
	NoteIdStr               string         `json:"note_id_str"`
	NoteLine                string         `json:"note_line"`
	NoteLineStyle           string         `json:"note_line_style"`
	NoteTitle               string         `json:"note_title"`
	NoteType                int            `json:"note_type"`
	NotesCount              NotesCount     `json:"notes_count"`
	NotesOwner              NotesUser      `json:"notes_owner"`
	NotesTopicInfo          NotesTopicInfo `json:"notes_topic_info,omitempty"`
	OriginAuditState        int            `json:"origin_audit_state"`
	OriginContentType       int            `json:"origin_content_type"`
	OriginNoteIdHazy        string         `json:"origin_note_id_hazy"`
	OriginNoteIdStr         string         `json:"origin_note_id_str"`
	OriginNotesOwner        NotesUser      `json:"origin_notes_owner"`
	OriginState             int            `json:"origin_state"`
	Pid                     int            `json:"pid"`
	PidStr                  string         `json:"pid_str"`
	Ptype                   int            `json:"ptype"`
	RefId                   int            `json:"ref_id"`
	RootAuditState          int            `json:"root_audit_state"`
	RootContentType         int            `json:"root_content_type"`
	RootHighlights          []interface{}  `json:"root_highlights"`
	RootNoteId              int            `json:"root_note_id"`
	RootNoteIdHazy          string         `json:"root_note_id_hazy"`
	RootNoteIdStr           string         `json:"root_note_id_str"`
	RootNotesOwner          NotesUser      `json:"root_notes_owner"`
	RootState               int            `json:"root_state"`
	ShareUrl                string         `json:"share_url"`
	SourceType              int            `json:"source_type"`
	State                   int            `json:"state"`
	StyleNoteLine           string         `json:"style_note_line"`
	Switch                  Switch         `json:"switch"`
	Tags                    []interface{}  `json:"tags"`
	Uid                     int            `json:"uid"`
	UidHazy                 string         `json:"uid_hazy"`
	UpdateTime              int            `json:"update_time"`
	UserExpectStatus        int            `json:"user_expect_status"`
	Video                   CommentVideo   `json:"video"`
	CommentIdStr            string         `json:"comment_id_str"`
	CommentReplyUser        NotesUser      `json:"comment_reply_user"`
	RepostCommentNoteIdStr  string         `json:"repost_comment_note_id_str"`
	RepostCommentNoteIdHazy string         `json:"repost_comment_note_id_hazy"`
	PidHazy                 string         `json:"pid_hazy"`
}
type Lesson struct {
	Pid     int    `json:"pid"`
	PidStr  string `json:"pid_str"`
	Ptype   int    `json:"ptype"`
	PidHazy string `json:"pid_hazy"`
}
type Switch struct {
	ImgOrigin bool `json:"img_origin"`
}

type CommentVideo struct {
	CardType             int    `json:"card_type"`
	Resource             string `json:"resource"`
	ResourceCommentCount int    `json:"resource_comment_count"`
	ResourceIcon         string `json:"resource_icon"`
	ResourceStudyCount   int    `json:"resource_study_count"`
	VideoCover           string `json:"video_cover"`
	VideoDuration        int    `json:"video_duration"`
	VideoDurationLabel   string `json:"video_duration_label"`
	VideoHeight          int    `json:"video_height"`
	VideoId              int    `json:"video_id"`
	VideoRst             string `json:"video_rst"`
	VideoState           int    `json:"video_state"`
	VideoWidth           int    `json:"video_width"`
	ViewCount            int    `json:"view_count"`
}

type Ddurl struct {
	NeedVisitorPopLoginView bool   `json:"needVisitorPopLoginView"`
	NeedCheckBuy            bool   `json:"needCheckBuy"`
	Url1                    string `json:"url1"`
	Url2                    string `json:"url2"`
}

type NotesCount struct {
	CommentCount int `json:"comment_count"`
	LikeCount    int `json:"like_count"`
	RepostCount  int `json:"repost_count"`
}

type NotesTopicInfo struct {
	IsTopicNotesElected bool   `json:"is_topic_notes_elected"`
	IsTopicNotesTopmost bool   `json:"is_topic_notes_topmost"`
	NotesTopicId        string `json:"notes_topic_id"`
	NotesTopicIdHazy    string `json:"notes_topic_id_hazy"`
	NotesTopicName      string `json:"notes_topic_name"`
}

type NotesExtra struct {
	OldClassID       int           `json:"OldClassID"`
	OldClassType     int           `json:"OldClassType"`
	ArticleTitle     string        `json:"article_title"`
	AuthorName       string        `json:"author_name"`
	BookAuthor       string        `json:"book_author"`
	BookId           int           `json:"book_id"`
	BookIsOldVersion int           `json:"book_is_old_version"`
	BookName         string        `json:"book_name"`
	BookShelfStatus  int           `json:"book_shelf_status"`
	ColumnTitle      string        `json:"column_title"`
	Images           []interface{} `json:"images"`
	ImagesSuffix     []interface{} `json:"images_suffix"`
	Img              string        `json:"img"`
	LogId            string        `json:"log_id"`
	ResourceIcon     string        `json:"resource_icon"`
	ScoreDesc        string        `json:"score_desc"`
	ScoreStr         string        `json:"score_str"`
	ShareUrl         string        `json:"share_url"`
	SourceId         int           `json:"source_id"`
	SourceSubTitle   string        `json:"source_sub_title"`
	SourceTitle      string        `json:"source_title"`
	SourceType       int           `json:"source_type"`
	SourceTypeName   string        `json:"source_type_name"`
	SubTitle         string        `json:"sub_title"`
	Title            string        `json:"title"`
	Tname            string        `json:"tname"`
	IsVipBook        int           `json:"is_vip_book,omitempty"`
	LogType          string        `json:"log_type,omitempty"`

	ColumnIntro  string `json:"column_intro"`
	ViewType     int    `json:"view_type"`
	AudioIdAlias string `json:"audio_id_alias"`
}
