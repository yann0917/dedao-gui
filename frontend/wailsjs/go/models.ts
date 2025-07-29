export namespace backend {
	
	export class LoginResult {
	    status: number;
	    cookie: string;
	    user?: services.User;
	
	    static createFrom(source: any = {}) {
	        return new LoginResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.cookie = source["cookie"];
	        this.user = this.convertValues(source["user"], services.User);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class QrCodeResp {
	    token: string;
	    qrCode: string;
	    qrCodeString: string;
	
	    static createFrom(source: any = {}) {
	        return new QrCodeResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.token = source["token"];
	        this.qrCode = source["qrCode"];
	        this.qrCodeString = source["qrCodeString"];
	    }
	}

}

export namespace options {
	
	export class SecondInstanceData {
	    Args: string[];
	    WorkingDirectory: string;
	
	    static createFrom(source: any = {}) {
	        return new SecondInstanceData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Args = source["Args"];
	        this.WorkingDirectory = source["WorkingDirectory"];
	    }
	}

}

export namespace services {
	
	export class AgencyDetail {
	    id: number;
	    id_str: string;
	    name: string;
	    intro: string;
	    qcg_id: number;
	    qcg_member_name: string;
	    qcg_member_avatar: string;
	    book_count: number;
	    uv: number;
	    status: number;
	
	    static createFrom(source: any = {}) {
	        return new AgencyDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.id_str = source["id_str"];
	        this.name = source["name"];
	        this.intro = source["intro"];
	        this.qcg_id = source["qcg_id"];
	        this.qcg_member_name = source["qcg_member_name"];
	        this.qcg_member_avatar = source["qcg_member_avatar"];
	        this.book_count = source["book_count"];
	        this.uv = source["uv"];
	        this.status = source["status"];
	    }
	}
	export class Option {
	    name: string;
	    value: string;
	    sub_options?: Option[];
	
	    static createFrom(source: any = {}) {
	        return new Option(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.value = source["value"];
	        this.sub_options = this.convertValues(source["sub_options"], Option);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Strategy {
	    title: string;
	    is_multiple: boolean;
	    is_hide: boolean;
	    options: Option[];
	
	    static createFrom(source: any = {}) {
	        return new Strategy(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.is_multiple = source["is_multiple"];
	        this.is_hide = source["is_hide"];
	        this.options = this.convertValues(source["options"], Option);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AlgoFilter {
	    title: string;
	    product_types: Strategy;
	    sort_strategy: Strategy;
	    progress_strategy: Strategy;
	    buy_strategy: Strategy;
	    navigations: Strategy;
	    tags: Strategy;
	
	    static createFrom(source: any = {}) {
	        return new AlgoFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.product_types = this.convertValues(source["product_types"], Strategy);
	        this.sort_strategy = this.convertValues(source["sort_strategy"], Strategy);
	        this.progress_strategy = this.convertValues(source["progress_strategy"], Strategy);
	        this.buy_strategy = this.convertValues(source["buy_strategy"], Strategy);
	        this.navigations = this.convertValues(source["navigations"], Strategy);
	        this.tags = this.convertValues(source["tags"], Strategy);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AlgoFilterParam {
	    classfc_name: string;
	    label_id: string;
	    nav_type: number;
	    navigation_id: string;
	    page: number;
	    page_size: number;
	    product_types: string;
	    request_id: string;
	    sort_strategy: string;
	    tags_ids: any[];
	
	    static createFrom(source: any = {}) {
	        return new AlgoFilterParam(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.classfc_name = source["classfc_name"];
	        this.label_id = source["label_id"];
	        this.nav_type = source["nav_type"];
	        this.navigation_id = source["navigation_id"];
	        this.page = source["page"];
	        this.page_size = source["page_size"];
	        this.product_types = source["product_types"];
	        this.request_id = source["request_id"];
	        this.sort_strategy = source["sort_strategy"];
	        this.tags_ids = source["tags_ids"];
	    }
	}
	export class AlgoFilterResp {
	    filter: AlgoFilter;
	    total: number;
	    request: AlgoFilterParam;
	
	    static createFrom(source: any = {}) {
	        return new AlgoFilterResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filter = this.convertValues(source["filter"], AlgoFilter);
	        this.total = source["total"];
	        this.request = this.convertValues(source["request"], AlgoFilterParam);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CostIntro {
	    price: string;
	
	    static createFrom(source: any = {}) {
	        return new CostIntro(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.price = source["price"];
	    }
	}
	export class ArticleSimpleInfo {
	    article_id: number;
	    origin_article_id: number;
	    article_audio_duration: number;
	    article_title: string;
	    article_intro: string;
	    image: string;
	
	    static createFrom(source: any = {}) {
	        return new ArticleSimpleInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.article_id = source["article_id"];
	        this.origin_article_id = source["origin_article_id"];
	        this.article_audio_duration = source["article_audio_duration"];
	        this.article_title = source["article_title"];
	        this.article_intro = source["article_intro"];
	        this.image = source["image"];
	    }
	}
	export class AlgoProduct {
	    item_type: number;
	    id: number;
	    product_type: number;
	    product_id: number;
	    name: string;
	    intro: string;
	    phase_num: number;
	    learn_user_count: number;
	    price: number;
	    current_price: number;
	    index_img: string;
	    horizontal_image: string;
	    lecturers_img: string;
	    price_desc: string;
	    is_buy: number;
	    is_vip: boolean;
	    trackinfo: string;
	    log_id: string;
	    log_type: string;
	    corner_img: string;
	    lecturer_name_and_title: string;
	    lecturers_v_status: number;
	    corner_img_vertical: string;
	    lecturer_name: string;
	    lecturer_title: string;
	    author_list: string[];
	    id_out: string;
	    is_limit_free: boolean;
	    has_play_auth: boolean;
	    button_type: number;
	    audio_id: string;
	    alias_id: string;
	    in_bookrack: boolean;
	    is_show_newbook: boolean;
	    is_subscribe: boolean;
	    online_time: string;
	    article_info: ArticleSimpleInfo;
	    need_login: number;
	    is_on_bookshelf: boolean;
	    duration: number;
	    collection_num: number;
	    dd_url: string;
	    learning_days_desc: string;
	    is_preferential: number;
	    is_count_down: number;
	    preferential_start_time: number;
	    preferential_end_time: number;
	    early_bird_price: number;
	    time_now: number;
	    hot_learn_desc: string;
	    score: string;
	    introduction: string;
	    user_score_count: number;
	    progress: number;
	    metrics: string;
	    cost_intro: CostIntro;
	    consumed_num: number;
	    free_maximum: number;
	    b_selling_channel_group: number;
	    track_info: string;
	    sort_values: any[];
	
	    static createFrom(source: any = {}) {
	        return new AlgoProduct(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.item_type = source["item_type"];
	        this.id = source["id"];
	        this.product_type = source["product_type"];
	        this.product_id = source["product_id"];
	        this.name = source["name"];
	        this.intro = source["intro"];
	        this.phase_num = source["phase_num"];
	        this.learn_user_count = source["learn_user_count"];
	        this.price = source["price"];
	        this.current_price = source["current_price"];
	        this.index_img = source["index_img"];
	        this.horizontal_image = source["horizontal_image"];
	        this.lecturers_img = source["lecturers_img"];
	        this.price_desc = source["price_desc"];
	        this.is_buy = source["is_buy"];
	        this.is_vip = source["is_vip"];
	        this.trackinfo = source["trackinfo"];
	        this.log_id = source["log_id"];
	        this.log_type = source["log_type"];
	        this.corner_img = source["corner_img"];
	        this.lecturer_name_and_title = source["lecturer_name_and_title"];
	        this.lecturers_v_status = source["lecturers_v_status"];
	        this.corner_img_vertical = source["corner_img_vertical"];
	        this.lecturer_name = source["lecturer_name"];
	        this.lecturer_title = source["lecturer_title"];
	        this.author_list = source["author_list"];
	        this.id_out = source["id_out"];
	        this.is_limit_free = source["is_limit_free"];
	        this.has_play_auth = source["has_play_auth"];
	        this.button_type = source["button_type"];
	        this.audio_id = source["audio_id"];
	        this.alias_id = source["alias_id"];
	        this.in_bookrack = source["in_bookrack"];
	        this.is_show_newbook = source["is_show_newbook"];
	        this.is_subscribe = source["is_subscribe"];
	        this.online_time = source["online_time"];
	        this.article_info = this.convertValues(source["article_info"], ArticleSimpleInfo);
	        this.need_login = source["need_login"];
	        this.is_on_bookshelf = source["is_on_bookshelf"];
	        this.duration = source["duration"];
	        this.collection_num = source["collection_num"];
	        this.dd_url = source["dd_url"];
	        this.learning_days_desc = source["learning_days_desc"];
	        this.is_preferential = source["is_preferential"];
	        this.is_count_down = source["is_count_down"];
	        this.preferential_start_time = source["preferential_start_time"];
	        this.preferential_end_time = source["preferential_end_time"];
	        this.early_bird_price = source["early_bird_price"];
	        this.time_now = source["time_now"];
	        this.hot_learn_desc = source["hot_learn_desc"];
	        this.score = source["score"];
	        this.introduction = source["introduction"];
	        this.user_score_count = source["user_score_count"];
	        this.progress = source["progress"];
	        this.metrics = source["metrics"];
	        this.cost_intro = this.convertValues(source["cost_intro"], CostIntro);
	        this.consumed_num = source["consumed_num"];
	        this.free_maximum = source["free_maximum"];
	        this.b_selling_channel_group = source["b_selling_channel_group"];
	        this.track_info = source["track_info"];
	        this.sort_values = source["sort_values"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AlgoProductResp {
	    product_list: AlgoProduct[];
	    request_id: string;
	    total: number;
	    is_more: number;
	
	    static createFrom(source: any = {}) {
	        return new AlgoProductResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.product_list = this.convertValues(source["product_list"], AlgoProduct);
	        this.request_id = source["request_id"];
	        this.total = source["total"];
	        this.is_more = source["is_more"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ArticleBase {
	    id: number;
	    id_str: string;
	    enid: string;
	    class_enid: string;
	    origin_id: number;
	    origin_id_str: string;
	    product_type: number;
	    product_id: number;
	    product_id_str: string;
	    class_id: number;
	    class_id_str: string;
	    chapter_id: number;
	    chapter_id_str: string;
	    title: string;
	    logo: string;
	    url: string;
	    summary: string;
	    mold: number;
	    push_content: string;
	    publish_time: number;
	    push_time: number;
	    push_status: number;
	    share_title: string;
	    share_content: string;
	    share_switch: number;
	    dd_article_id: number;
	    dd_article_id_str: string;
	    dd_article_token: string;
	    status: number;
	    create_time: number;
	    update_time: number;
	    cur_learn_count: number;
	    is_free_try: boolean;
	    is_user_free_try: boolean;
	    order_num: number;
	    is_like: boolean;
	    share_url: string;
	    trial_share_url: string;
	    is_read: boolean;
	    log_id: string;
	    log_type: string;
	    recommend_title: string;
	    audio_alias_ids: string[];
	    is_buy: boolean;
	    dd_media_id: number;
	    dd_media_id_str: string;
	    video_status: number;
	    dd_live_id: number;
	    not_join_plan: number;
	
	    static createFrom(source: any = {}) {
	        return new ArticleBase(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.id_str = source["id_str"];
	        this.enid = source["enid"];
	        this.class_enid = source["class_enid"];
	        this.origin_id = source["origin_id"];
	        this.origin_id_str = source["origin_id_str"];
	        this.product_type = source["product_type"];
	        this.product_id = source["product_id"];
	        this.product_id_str = source["product_id_str"];
	        this.class_id = source["class_id"];
	        this.class_id_str = source["class_id_str"];
	        this.chapter_id = source["chapter_id"];
	        this.chapter_id_str = source["chapter_id_str"];
	        this.title = source["title"];
	        this.logo = source["logo"];
	        this.url = source["url"];
	        this.summary = source["summary"];
	        this.mold = source["mold"];
	        this.push_content = source["push_content"];
	        this.publish_time = source["publish_time"];
	        this.push_time = source["push_time"];
	        this.push_status = source["push_status"];
	        this.share_title = source["share_title"];
	        this.share_content = source["share_content"];
	        this.share_switch = source["share_switch"];
	        this.dd_article_id = source["dd_article_id"];
	        this.dd_article_id_str = source["dd_article_id_str"];
	        this.dd_article_token = source["dd_article_token"];
	        this.status = source["status"];
	        this.create_time = source["create_time"];
	        this.update_time = source["update_time"];
	        this.cur_learn_count = source["cur_learn_count"];
	        this.is_free_try = source["is_free_try"];
	        this.is_user_free_try = source["is_user_free_try"];
	        this.order_num = source["order_num"];
	        this.is_like = source["is_like"];
	        this.share_url = source["share_url"];
	        this.trial_share_url = source["trial_share_url"];
	        this.is_read = source["is_read"];
	        this.log_id = source["log_id"];
	        this.log_type = source["log_type"];
	        this.recommend_title = source["recommend_title"];
	        this.audio_alias_ids = source["audio_alias_ids"];
	        this.is_buy = source["is_buy"];
	        this.dd_media_id = source["dd_media_id"];
	        this.dd_media_id_str = source["dd_media_id_str"];
	        this.video_status = source["video_status"];
	        this.dd_live_id = source["dd_live_id"];
	        this.not_join_plan = source["not_join_plan"];
	    }
	}
	export class Video {
	    token: string;
	    token_version: number;
	    cover_img: string;
	    dd_media_id: number;
	    dd_media_id_str: string;
	    duration: number;
	    bitrate_480: string;
	    bitrate_480_size: number;
	    bitrate_480_audio: string;
	    bitrate_720: string;
	    bitrate_720_size: number;
	    bitrate_720_audio: string;
	    bitrate_1080: string;
	    bitrate_1080_size: number;
	    bitrate_1080_audio: string;
	    is_drm: boolean;
	    listen_progress: number;
	    listen_finished: boolean;
	    log_id: string;
	    log_type: string;
	    caption: string;
	
	    static createFrom(source: any = {}) {
	        return new Video(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.token = source["token"];
	        this.token_version = source["token_version"];
	        this.cover_img = source["cover_img"];
	        this.dd_media_id = source["dd_media_id"];
	        this.dd_media_id_str = source["dd_media_id_str"];
	        this.duration = source["duration"];
	        this.bitrate_480 = source["bitrate_480"];
	        this.bitrate_480_size = source["bitrate_480_size"];
	        this.bitrate_480_audio = source["bitrate_480_audio"];
	        this.bitrate_720 = source["bitrate_720"];
	        this.bitrate_720_size = source["bitrate_720_size"];
	        this.bitrate_720_audio = source["bitrate_720_audio"];
	        this.bitrate_1080 = source["bitrate_1080"];
	        this.bitrate_1080_size = source["bitrate_1080_size"];
	        this.bitrate_1080_audio = source["bitrate_1080_audio"];
	        this.is_drm = source["is_drm"];
	        this.listen_progress = source["listen_progress"];
	        this.listen_finished = source["listen_finished"];
	        this.log_id = source["log_id"];
	        this.log_type = source["log_type"];
	        this.caption = source["caption"];
	    }
	}
	export class Audio {
	    alias_id: string;
	    extrainfo: string;
	    class_id: number;
	    title: string;
	    share_title: string;
	    mp3_play_url: string;
	    duration: number;
	    schedule: number;
	    topic_id: number;
	    summary: string;
	    price: number;
	    icon: string;
	    size: number;
	    etag: string;
	    token: string;
	    share_summary: string;
	    collection: number;
	    count: number;
	    bored_count: number;
	    audio_type: number;
	    drm_version: number;
	    source_id: number;
	    source_type: number;
	    source_icon: string;
	    source_name: string;
	    listen_progress: number;
	    listen_finished: boolean;
	    dd_article_id: string;
	    audio_list_icon: string;
	    class_course_id: string;
	    class_article_id: string;
	    log_id: string;
	    log_type: string;
	    log_interface: string;
	    trackinfo: string;
	    used_drm: number;
	    index_img: string;
	    reader: string;
	    reader_name: string;
	    odob_group_enid: string;
	
	    static createFrom(source: any = {}) {
	        return new Audio(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.alias_id = source["alias_id"];
	        this.extrainfo = source["extrainfo"];
	        this.class_id = source["class_id"];
	        this.title = source["title"];
	        this.share_title = source["share_title"];
	        this.mp3_play_url = source["mp3_play_url"];
	        this.duration = source["duration"];
	        this.schedule = source["schedule"];
	        this.topic_id = source["topic_id"];
	        this.summary = source["summary"];
	        this.price = source["price"];
	        this.icon = source["icon"];
	        this.size = source["size"];
	        this.etag = source["etag"];
	        this.token = source["token"];
	        this.share_summary = source["share_summary"];
	        this.collection = source["collection"];
	        this.count = source["count"];
	        this.bored_count = source["bored_count"];
	        this.audio_type = source["audio_type"];
	        this.drm_version = source["drm_version"];
	        this.source_id = source["source_id"];
	        this.source_type = source["source_type"];
	        this.source_icon = source["source_icon"];
	        this.source_name = source["source_name"];
	        this.listen_progress = source["listen_progress"];
	        this.listen_finished = source["listen_finished"];
	        this.dd_article_id = source["dd_article_id"];
	        this.audio_list_icon = source["audio_list_icon"];
	        this.class_course_id = source["class_course_id"];
	        this.class_article_id = source["class_article_id"];
	        this.log_id = source["log_id"];
	        this.log_type = source["log_type"];
	        this.log_interface = source["log_interface"];
	        this.trackinfo = source["trackinfo"];
	        this.used_drm = source["used_drm"];
	        this.index_img = source["index_img"];
	        this.reader = source["reader"];
	        this.reader_name = source["reader_name"];
	        this.odob_group_enid = source["odob_group_enid"];
	    }
	}
	export class MediaBaseInfo {
	    media_type: number;
	    source_id: string;
	    security_token: string;
	
	    static createFrom(source: any = {}) {
	        return new MediaBaseInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.media_type = source["media_type"];
	        this.source_id = source["source_id"];
	        this.security_token = source["security_token"];
	    }
	}
	export class ArticleIntro {
	    id: number;
	    id_str: string;
	    enid: string;
	    class_enid: string;
	    origin_id: number;
	    origin_id_str: string;
	    product_type: number;
	    product_id: number;
	    product_id_str: string;
	    class_id: number;
	    class_id_str: string;
	    chapter_id: number;
	    chapter_id_str: string;
	    title: string;
	    logo: string;
	    url: string;
	    summary: string;
	    mold: number;
	    push_content: string;
	    publish_time: number;
	    push_time: number;
	    push_status: number;
	    share_title: string;
	    share_content: string;
	    share_switch: number;
	    dd_article_id: number;
	    dd_article_id_str: string;
	    dd_article_token: string;
	    status: number;
	    create_time: number;
	    update_time: number;
	    cur_learn_count: number;
	    is_free_try: boolean;
	    is_user_free_try: boolean;
	    order_num: number;
	    is_like: boolean;
	    share_url: string;
	    trial_share_url: string;
	    is_read: boolean;
	    log_id: string;
	    log_type: string;
	    recommend_title: string;
	    audio_alias_ids: string[];
	    is_buy: boolean;
	    dd_media_id: number;
	    dd_media_id_str: string;
	    video_status: number;
	    dd_live_id: number;
	    not_join_plan: number;
	    media_base_info: MediaBaseInfo[];
	    audio?: Audio;
	    video?: Video[];
	
	    static createFrom(source: any = {}) {
	        return new ArticleIntro(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.id_str = source["id_str"];
	        this.enid = source["enid"];
	        this.class_enid = source["class_enid"];
	        this.origin_id = source["origin_id"];
	        this.origin_id_str = source["origin_id_str"];
	        this.product_type = source["product_type"];
	        this.product_id = source["product_id"];
	        this.product_id_str = source["product_id_str"];
	        this.class_id = source["class_id"];
	        this.class_id_str = source["class_id_str"];
	        this.chapter_id = source["chapter_id"];
	        this.chapter_id_str = source["chapter_id_str"];
	        this.title = source["title"];
	        this.logo = source["logo"];
	        this.url = source["url"];
	        this.summary = source["summary"];
	        this.mold = source["mold"];
	        this.push_content = source["push_content"];
	        this.publish_time = source["publish_time"];
	        this.push_time = source["push_time"];
	        this.push_status = source["push_status"];
	        this.share_title = source["share_title"];
	        this.share_content = source["share_content"];
	        this.share_switch = source["share_switch"];
	        this.dd_article_id = source["dd_article_id"];
	        this.dd_article_id_str = source["dd_article_id_str"];
	        this.dd_article_token = source["dd_article_token"];
	        this.status = source["status"];
	        this.create_time = source["create_time"];
	        this.update_time = source["update_time"];
	        this.cur_learn_count = source["cur_learn_count"];
	        this.is_free_try = source["is_free_try"];
	        this.is_user_free_try = source["is_user_free_try"];
	        this.order_num = source["order_num"];
	        this.is_like = source["is_like"];
	        this.share_url = source["share_url"];
	        this.trial_share_url = source["trial_share_url"];
	        this.is_read = source["is_read"];
	        this.log_id = source["log_id"];
	        this.log_type = source["log_type"];
	        this.recommend_title = source["recommend_title"];
	        this.audio_alias_ids = source["audio_alias_ids"];
	        this.is_buy = source["is_buy"];
	        this.dd_media_id = source["dd_media_id"];
	        this.dd_media_id_str = source["dd_media_id_str"];
	        this.video_status = source["video_status"];
	        this.dd_live_id = source["dd_live_id"];
	        this.not_join_plan = source["not_join_plan"];
	        this.media_base_info = this.convertValues(source["media_base_info"], MediaBaseInfo);
	        this.audio = this.convertValues(source["audio"], Audio);
	        this.video = this.convertValues(source["video"], Video);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ArticleList {
	    article_list: ArticleIntro[];
	    class_id: number;
	    ptype: number;
	    pid: number;
	    reverse: boolean;
	    chapter_id: string;
	    max_id: number;
	
	    static createFrom(source: any = {}) {
	        return new ArticleList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.article_list = this.convertValues(source["article_list"], ArticleIntro);
	        this.class_id = source["class_id"];
	        this.ptype = source["ptype"];
	        this.pid = source["pid"];
	        this.reverse = source["reverse"];
	        this.chapter_id = source["chapter_id"];
	        this.max_id = source["max_id"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class Rank {
	    rank_number: number;
	    rank_name: string;
	    rank_type: number;
	    rank_desc: string;
	    dd_url: string;
	
	    static createFrom(source: any = {}) {
	        return new Rank(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rank_number = source["rank_number"];
	        this.rank_name = source["rank_name"];
	        this.rank_type = source["rank_type"];
	        this.rank_desc = source["rank_desc"];
	        this.dd_url = source["dd_url"];
	    }
	}
	export class LimitFree {
	    limit_free_expire: boolean;
	    is_limit_free: boolean;
	    is_red_packet: boolean;
	    free_begin_time: number;
	    free_end_time: number;
	    free_maximum: number;
	    consume_num: number;
	    vip_end_time: number;
	
	    static createFrom(source: any = {}) {
	        return new LimitFree(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.limit_free_expire = source["limit_free_expire"];
	        this.is_limit_free = source["is_limit_free"];
	        this.is_red_packet = source["is_red_packet"];
	        this.free_begin_time = source["free_begin_time"];
	        this.free_end_time = source["free_end_time"];
	        this.free_maximum = source["free_maximum"];
	        this.consume_num = source["consume_num"];
	        this.vip_end_time = source["vip_end_time"];
	    }
	}
	export class TopicSummary {
	    id: number;
	    topic_id: number;
	    title: string;
	    sub_title: string;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new TopicSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.topic_id = source["topic_id"];
	        this.title = source["title"];
	        this.sub_title = source["sub_title"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}
	export class AudioInfo {
	    id: number;
	    audio_id: string;
	    title: string;
	    icon: string;
	    duration: number;
	    audio_price: string;
	    audio_summary: string;
	    publish_time: number;
	    is_vip_gived: boolean;
	    is_vip: boolean;
	    agency_detail: AgencyDetail;
	    log_id: number;
	    is_buy: boolean;
	    in_bookrack: boolean;
	    progress: number;
	    tag: string[];
	    topic_summary: TopicSummary[];
	    is_limit_free: boolean;
	    topic_encode_id: string;
	    h5_share_url: string;
	    can_share: boolean;
	    limit_free: LimitFree;
	    vip_end_time: number;
	    button_type: number;
	    has_play_auth: boolean;
	    rank: Rank;
	    learn_count: number;
	    learn_count_desc: string;
	
	    static createFrom(source: any = {}) {
	        return new AudioInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.audio_id = source["audio_id"];
	        this.title = source["title"];
	        this.icon = source["icon"];
	        this.duration = source["duration"];
	        this.audio_price = source["audio_price"];
	        this.audio_summary = source["audio_summary"];
	        this.publish_time = source["publish_time"];
	        this.is_vip_gived = source["is_vip_gived"];
	        this.is_vip = source["is_vip"];
	        this.agency_detail = this.convertValues(source["agency_detail"], AgencyDetail);
	        this.log_id = source["log_id"];
	        this.is_buy = source["is_buy"];
	        this.in_bookrack = source["in_bookrack"];
	        this.progress = source["progress"];
	        this.tag = source["tag"];
	        this.topic_summary = this.convertValues(source["topic_summary"], TopicSummary);
	        this.is_limit_free = source["is_limit_free"];
	        this.topic_encode_id = source["topic_encode_id"];
	        this.h5_share_url = source["h5_share_url"];
	        this.can_share = source["can_share"];
	        this.limit_free = this.convertValues(source["limit_free"], LimitFree);
	        this.vip_end_time = source["vip_end_time"];
	        this.button_type = source["button_type"];
	        this.has_play_auth = source["has_play_auth"];
	        this.rank = this.convertValues(source["rank"], Rank);
	        this.learn_count = source["learn_count"];
	        this.learn_count_desc = source["learn_count_desc"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class QualityUser {
	    avatar: string;
	    id: string;
	    name: string;
	    title: string;
	
	    static createFrom(source: any = {}) {
	        return new QualityUser(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.avatar = source["avatar"];
	        this.id = source["id"];
	        this.name = source["name"];
	        this.title = source["title"];
	    }
	}
	export class Quality {
	    desc: string;
	    list: QualityUser[];
	    paytime: string;
	
	    static createFrom(source: any = {}) {
	        return new Quality(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.desc = source["desc"];
	        this.list = this.convertValues(source["list"], QualityUser);
	        this.paytime = source["paytime"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AudioInfoResp {
	    detail: AudioInfo;
	    quality: Quality;
	
	    static createFrom(source: any = {}) {
	        return new AudioInfoResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.detail = this.convertValues(source["detail"], AudioInfo);
	        this.quality = this.convertValues(source["quality"], Quality);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Banner {
	    beginTime: number;
	    endTime: number;
	    id: number;
	    img: string;
	    isDelete: number;
	    moduleId: number;
	    sort: number;
	    sourceId: string;
	    title: string;
	    url: string;
	    urlType: number;
	
	    static createFrom(source: any = {}) {
	        return new Banner(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.beginTime = source["beginTime"];
	        this.endTime = source["endTime"];
	        this.id = source["id"];
	        this.img = source["img"];
	        this.isDelete = source["isDelete"];
	        this.moduleId = source["moduleId"];
	        this.sort = source["sort"];
	        this.sourceId = source["sourceId"];
	        this.title = source["title"];
	        this.url = source["url"];
	        this.urlType = source["urlType"];
	    }
	}
	export class Book {
	    book_type: number;
	
	    static createFrom(source: any = {}) {
	        return new Book(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.book_type = source["book_type"];
	    }
	}
	export class CampHazy {
	    SourceIdHazy: string;
	    OriginIdHazy: string;
	
	    static createFrom(source: any = {}) {
	        return new CampHazy(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.SourceIdHazy = source["SourceIdHazy"];
	        this.OriginIdHazy = source["OriginIdHazy"];
	    }
	}
	export class Catalog {
	    level: number;
	    text: string;
	    href: string;
	    playOrder: number;
	
	    static createFrom(source: any = {}) {
	        return new Catalog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.level = source["level"];
	        this.text = source["text"];
	        this.href = source["href"];
	        this.playOrder = source["playOrder"];
	    }
	}
	export class Chapter {
	    id: number;
	    id_str: string;
	    class_id: number;
	    class_id_str: string;
	    order_name: string;
	    name: string;
	    intro: string;
	    phase_num: number;
	    status: number;
	    order_num: number;
	    is_finished: number;
	    update_time: number;
	    log_id: string;
	    log_type: string;
	    appendage_count: number;
	    article_list: ArticleBase[];
	
	    static createFrom(source: any = {}) {
	        return new Chapter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.id_str = source["id_str"];
	        this.class_id = source["class_id"];
	        this.class_id_str = source["class_id_str"];
	        this.order_name = source["order_name"];
	        this.name = source["name"];
	        this.intro = source["intro"];
	        this.phase_num = source["phase_num"];
	        this.status = source["status"];
	        this.order_num = source["order_num"];
	        this.is_finished = source["is_finished"];
	        this.update_time = source["update_time"];
	        this.log_id = source["log_id"];
	        this.log_type = source["log_type"];
	        this.appendage_count = source["appendage_count"];
	        this.article_list = this.convertValues(source["article_list"], ArticleBase);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ClassComment {
	    id: number;
	    user_id: number;
	    title: string;
	    score: number;
	    content: string;
	    no_style_content: string;
	    note_id: string;
	    nickname: string;
	    avatar: string;
	    avatar_s: string;
	    track_info: string;
	
	    static createFrom(source: any = {}) {
	        return new ClassComment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.user_id = source["user_id"];
	        this.title = source["title"];
	        this.score = source["score"];
	        this.content = source["content"];
	        this.no_style_content = source["no_style_content"];
	        this.note_id = source["note_id"];
	        this.nickname = source["nickname"];
	        this.avatar = source["avatar"];
	        this.avatar_s = source["avatar_s"];
	        this.track_info = source["track_info"];
	    }
	}
	export class ClassCommentInfo {
	    comment_list: ClassComment[];
	    count: number;
	    average_score: string;
	
	    static createFrom(source: any = {}) {
	        return new ClassCommentInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.comment_list = this.convertValues(source["comment_list"], ClassComment);
	        this.count = source["count"];
	        this.average_score = source["average_score"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Collection {
	    is_collected: boolean;
	    collection_count: number;
	
	    static createFrom(source: any = {}) {
	        return new Collection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.is_collected = source["is_collected"];
	        this.collection_count = source["collection_count"];
	    }
	}
	export class ClassInfo {
	    log_id: string;
	    log_type: string;
	    id: number;
	    id_str: string;
	    enid: string;
	    product_id: number;
	    product_type: number;
	    has_chapter: number;
	    name: string;
	    intro: string;
	    phase_num: number;
	    learn_user_count: number;
	    current_article_count: number;
	    highlight: string;
	    price: number;
	    index_img: string;
	    index_img_applet: string;
	    logo: string;
	    logo_iphonex: string;
	    square_img: string;
	    outline_img: string;
	    player_img: string;
	    share_title: string;
	    share_summary: string;
	    status: number;
	    order_num: number;
	    shzf_url: string;
	    shzf_url_qr: string;
	    price_desc: string;
	    article_time: number;
	    article_title: string;
	    is_subscribe: number;
	    lecturer_uid: number;
	    lecturer_name: string;
	    lecturer_title: string;
	    lecturer_intro: string;
	    lecturer_name_and_title: string;
	    lecturer_avatar: string;
	    is_finished: number;
	    update_time: number;
	    share_url: string;
	    default_sort_reverse: boolean;
	    presale_url: string;
	    without_audio: boolean;
	    view_type: number;
	    h5_url_name: string;
	    package_manager_switch: boolean;
	    package_manager: boolean;
	    lecture_avator_special: string;
	    mini_share_img: string;
	    estimated_shelf_time: number;
	    estimated_down_time: number;
	    corner_img: string;
	    corner_img_vertical: string;
	    without_giving: boolean;
	    dd_url: string;
	    publish_time: number;
	    dd_media_id: string;
	    video_cover: string;
	    is_in_vip: boolean;
	    is_vip: boolean;
	    collection: Collection;
	    formal_article_count: number;
	    video_class: number;
	    video_class_intro: string;
	    activity_icon: string;
	    activity_title: string;
	    activity_url: string;
	    reality_extra_count: number;
	    reality_formal_count: number;
	    intro_article_ids: number[];
	    is_preferential: number;
	    is_count_down: number;
	    preferential_start_time: number;
	    preferential_end_time: number;
	    early_bird_price: number;
	    trial_count: number;
	    special_logo: string;
	    early_bird_msg: string;
	
	    static createFrom(source: any = {}) {
	        return new ClassInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.log_id = source["log_id"];
	        this.log_type = source["log_type"];
	        this.id = source["id"];
	        this.id_str = source["id_str"];
	        this.enid = source["enid"];
	        this.product_id = source["product_id"];
	        this.product_type = source["product_type"];
	        this.has_chapter = source["has_chapter"];
	        this.name = source["name"];
	        this.intro = source["intro"];
	        this.phase_num = source["phase_num"];
	        this.learn_user_count = source["learn_user_count"];
	        this.current_article_count = source["current_article_count"];
	        this.highlight = source["highlight"];
	        this.price = source["price"];
	        this.index_img = source["index_img"];
	        this.index_img_applet = source["index_img_applet"];
	        this.logo = source["logo"];
	        this.logo_iphonex = source["logo_iphonex"];
	        this.square_img = source["square_img"];
	        this.outline_img = source["outline_img"];
	        this.player_img = source["player_img"];
	        this.share_title = source["share_title"];
	        this.share_summary = source["share_summary"];
	        this.status = source["status"];
	        this.order_num = source["order_num"];
	        this.shzf_url = source["shzf_url"];
	        this.shzf_url_qr = source["shzf_url_qr"];
	        this.price_desc = source["price_desc"];
	        this.article_time = source["article_time"];
	        this.article_title = source["article_title"];
	        this.is_subscribe = source["is_subscribe"];
	        this.lecturer_uid = source["lecturer_uid"];
	        this.lecturer_name = source["lecturer_name"];
	        this.lecturer_title = source["lecturer_title"];
	        this.lecturer_intro = source["lecturer_intro"];
	        this.lecturer_name_and_title = source["lecturer_name_and_title"];
	        this.lecturer_avatar = source["lecturer_avatar"];
	        this.is_finished = source["is_finished"];
	        this.update_time = source["update_time"];
	        this.share_url = source["share_url"];
	        this.default_sort_reverse = source["default_sort_reverse"];
	        this.presale_url = source["presale_url"];
	        this.without_audio = source["without_audio"];
	        this.view_type = source["view_type"];
	        this.h5_url_name = source["h5_url_name"];
	        this.package_manager_switch = source["package_manager_switch"];
	        this.package_manager = source["package_manager"];
	        this.lecture_avator_special = source["lecture_avator_special"];
	        this.mini_share_img = source["mini_share_img"];
	        this.estimated_shelf_time = source["estimated_shelf_time"];
	        this.estimated_down_time = source["estimated_down_time"];
	        this.corner_img = source["corner_img"];
	        this.corner_img_vertical = source["corner_img_vertical"];
	        this.without_giving = source["without_giving"];
	        this.dd_url = source["dd_url"];
	        this.publish_time = source["publish_time"];
	        this.dd_media_id = source["dd_media_id"];
	        this.video_cover = source["video_cover"];
	        this.is_in_vip = source["is_in_vip"];
	        this.is_vip = source["is_vip"];
	        this.collection = this.convertValues(source["collection"], Collection);
	        this.formal_article_count = source["formal_article_count"];
	        this.video_class = source["video_class"];
	        this.video_class_intro = source["video_class_intro"];
	        this.activity_icon = source["activity_icon"];
	        this.activity_title = source["activity_title"];
	        this.activity_url = source["activity_url"];
	        this.reality_extra_count = source["reality_extra_count"];
	        this.reality_formal_count = source["reality_formal_count"];
	        this.intro_article_ids = source["intro_article_ids"];
	        this.is_preferential = source["is_preferential"];
	        this.is_count_down = source["is_count_down"];
	        this.preferential_start_time = source["preferential_start_time"];
	        this.preferential_end_time = source["preferential_end_time"];
	        this.early_bird_price = source["early_bird_price"];
	        this.trial_count = source["trial_count"];
	        this.special_logo = source["special_logo"];
	        this.early_bird_msg = source["early_bird_msg"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ClassReviews {
	    show_text: string;
	    review_done: boolean;
	    is_satisfied: boolean;
	    detail: string;
	
	    static createFrom(source: any = {}) {
	        return new ClassReviews(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.show_text = source["show_text"];
	        this.review_done = source["review_done"];
	        this.is_satisfied = source["is_satisfied"];
	        this.detail = source["detail"];
	    }
	}
	
	export class Comb {
	    uid: number;
	    uid_hazy: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Comb(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uid = source["uid"];
	        this.uid_hazy = source["uid_hazy"];
	        this.name = source["name"];
	    }
	}
	export class CommentVideo {
	    card_type: number;
	    resource: string;
	    resource_comment_count: number;
	    resource_icon: string;
	    resource_study_count: number;
	    video_cover: string;
	    video_duration: number;
	    video_duration_label: string;
	    video_height: number;
	    video_id: number;
	    video_rst: string;
	    video_state: number;
	    video_width: number;
	    view_count: number;
	
	    static createFrom(source: any = {}) {
	        return new CommentVideo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.card_type = source["card_type"];
	        this.resource = source["resource"];
	        this.resource_comment_count = source["resource_comment_count"];
	        this.resource_icon = source["resource_icon"];
	        this.resource_study_count = source["resource_study_count"];
	        this.video_cover = source["video_cover"];
	        this.video_duration = source["video_duration"];
	        this.video_duration_label = source["video_duration_label"];
	        this.video_height = source["video_height"];
	        this.video_id = source["video_id"];
	        this.video_rst = source["video_rst"];
	        this.video_state = source["video_state"];
	        this.video_width = source["video_width"];
	        this.view_count = source["view_count"];
	    }
	}
	export class Switch {
	    img_origin: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Switch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.img_origin = source["img_origin"];
	    }
	}
	export class NotesTopicInfo {
	    is_topic_notes_elected: boolean;
	    is_topic_notes_topmost: boolean;
	    notes_topic_id: string;
	    notes_topic_id_hazy: string;
	    notes_topic_name: string;
	
	    static createFrom(source: any = {}) {
	        return new NotesTopicInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.is_topic_notes_elected = source["is_topic_notes_elected"];
	        this.is_topic_notes_topmost = source["is_topic_notes_topmost"];
	        this.notes_topic_id = source["notes_topic_id"];
	        this.notes_topic_id_hazy = source["notes_topic_id_hazy"];
	        this.notes_topic_name = source["notes_topic_name"];
	    }
	}
	export class NotesUser {
	    id: string;
	    uid: number;
	    name: string;
	    avatar: string;
	    follow: number;
	    isV: number;
	    slogan: string;
	    Vinfo: string;
	    student_id: number;
	    is_poster: boolean;
	    qrcode: string;
	    log_id: string;
	    log_type: string;
	    uid_hazy: string;
	    note_id_hazy: string;
	    role: string;
	    attribution: string;
	
	    static createFrom(source: any = {}) {
	        return new NotesUser(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.uid = source["uid"];
	        this.name = source["name"];
	        this.avatar = source["avatar"];
	        this.follow = source["follow"];
	        this.isV = source["isV"];
	        this.slogan = source["slogan"];
	        this.Vinfo = source["Vinfo"];
	        this.student_id = source["student_id"];
	        this.is_poster = source["is_poster"];
	        this.qrcode = source["qrcode"];
	        this.log_id = source["log_id"];
	        this.log_type = source["log_type"];
	        this.uid_hazy = source["uid_hazy"];
	        this.note_id_hazy = source["note_id_hazy"];
	        this.role = source["role"];
	        this.attribution = source["attribution"];
	    }
	}
	export class NotesCount {
	    comment_count: number;
	    like_count: number;
	    repost_count: number;
	
	    static createFrom(source: any = {}) {
	        return new NotesCount(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.comment_count = source["comment_count"];
	        this.like_count = source["like_count"];
	        this.repost_count = source["repost_count"];
	    }
	}
	export class Lesson {
	    pid: number;
	    pid_str: string;
	    ptype: number;
	    pid_hazy: string;
	
	    static createFrom(source: any = {}) {
	        return new Lesson(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pid = source["pid"];
	        this.pid_str = source["pid_str"];
	        this.ptype = source["ptype"];
	        this.pid_hazy = source["pid_hazy"];
	    }
	}
	export class NotesExtra {
	    OldClassID: number;
	    OldClassType: number;
	    article_title: string;
	    author_name: string;
	    book_author: string;
	    book_id: number;
	    book_is_old_version: number;
	    book_name: string;
	    book_shelf_status: number;
	    column_title: string;
	    images: any[];
	    images_suffix: any[];
	    img: string;
	    log_id: string;
	    resource_icon: string;
	    score_desc: string;
	    score_str: string;
	    share_url: string;
	    source_id: number;
	    source_sub_title: string;
	    source_title: string;
	    source_type: number;
	    source_type_name: string;
	    sub_title: string;
	    title: string;
	    tname: string;
	    is_vip_book?: number;
	    log_type?: string;
	    column_intro: string;
	    view_type: number;
	    audio_id_alias: string;
	
	    static createFrom(source: any = {}) {
	        return new NotesExtra(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.OldClassID = source["OldClassID"];
	        this.OldClassType = source["OldClassType"];
	        this.article_title = source["article_title"];
	        this.author_name = source["author_name"];
	        this.book_author = source["book_author"];
	        this.book_id = source["book_id"];
	        this.book_is_old_version = source["book_is_old_version"];
	        this.book_name = source["book_name"];
	        this.book_shelf_status = source["book_shelf_status"];
	        this.column_title = source["column_title"];
	        this.images = source["images"];
	        this.images_suffix = source["images_suffix"];
	        this.img = source["img"];
	        this.log_id = source["log_id"];
	        this.resource_icon = source["resource_icon"];
	        this.score_desc = source["score_desc"];
	        this.score_str = source["score_str"];
	        this.share_url = source["share_url"];
	        this.source_id = source["source_id"];
	        this.source_sub_title = source["source_sub_title"];
	        this.source_title = source["source_title"];
	        this.source_type = source["source_type"];
	        this.source_type_name = source["source_type_name"];
	        this.sub_title = source["sub_title"];
	        this.title = source["title"];
	        this.tname = source["tname"];
	        this.is_vip_book = source["is_vip_book"];
	        this.log_type = source["log_type"];
	        this.column_intro = source["column_intro"];
	        this.view_type = source["view_type"];
	        this.audio_id_alias = source["audio_id_alias"];
	    }
	}
	export class Ddurl {
	    needVisitorPopLoginView: boolean;
	    needCheckBuy: boolean;
	    url1: string;
	    url2: string;
	
	    static createFrom(source: any = {}) {
	        return new Ddurl(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.needVisitorPopLoginView = source["needVisitorPopLoginView"];
	        this.needCheckBuy = source["needCheckBuy"];
	        this.url1 = source["url1"];
	        this.url2 = source["url2"];
	    }
	}
	export class Comment {
	    attachment_type: number;
	    audit_state: number;
	    class: number;
	    comment_reply: string;
	    comment_reply_time: number;
	    content: string;
	    content_type: number;
	    create_time: number;
	    current_state: number;
	    ddurl: Ddurl;
	    detail_id: number;
	    extra: NotesExtra;
	    feed_id: number;
	    highlights: any[];
	    is_from_me: number;
	    can_edit: boolean;
	    is_like: boolean;
	    is_open_ledgers: boolean;
	    is_permission: boolean;
	    is_reposted: boolean;
	    is_upmost: boolean;
	    level_permission: boolean;
	    lesson: Lesson;
	    level: number;
	    level_type: number;
	    log_id: string;
	    log_type: string;
	    note: string;
	    note_id: number;
	    note_id_hazy: string;
	    note_id_str: string;
	    note_line: string;
	    note_line_style: string;
	    note_title: string;
	    note_type: number;
	    notes_count: NotesCount;
	    notes_owner: NotesUser;
	    notes_topic_info?: NotesTopicInfo;
	    origin_audit_state: number;
	    origin_content_type: number;
	    origin_note_id_hazy: string;
	    origin_note_id_str: string;
	    origin_notes_owner: NotesUser;
	    origin_state: number;
	    pid: number;
	    pid_str: string;
	    ptype: number;
	    ref_id: number;
	    root_audit_state: number;
	    root_content_type: number;
	    root_highlights: any[];
	    root_note_id: number;
	    root_note_id_hazy: string;
	    root_note_id_str: string;
	    root_notes_owner: NotesUser;
	    root_state: number;
	    share_url: string;
	    source_type: number;
	    state: number;
	    style_note_line: string;
	    switch: Switch;
	    tags: any[];
	    uid: number;
	    uid_hazy: string;
	    update_time: number;
	    user_expect_status: number;
	    video: CommentVideo;
	    comment_id_str: string;
	    comment_reply_user: NotesUser;
	    repost_comment_note_id_str: string;
	    repost_comment_note_id_hazy: string;
	    pid_hazy: string;
	
	    static createFrom(source: any = {}) {
	        return new Comment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.attachment_type = source["attachment_type"];
	        this.audit_state = source["audit_state"];
	        this.class = source["class"];
	        this.comment_reply = source["comment_reply"];
	        this.comment_reply_time = source["comment_reply_time"];
	        this.content = source["content"];
	        this.content_type = source["content_type"];
	        this.create_time = source["create_time"];
	        this.current_state = source["current_state"];
	        this.ddurl = this.convertValues(source["ddurl"], Ddurl);
	        this.detail_id = source["detail_id"];
	        this.extra = this.convertValues(source["extra"], NotesExtra);
	        this.feed_id = source["feed_id"];
	        this.highlights = source["highlights"];
	        this.is_from_me = source["is_from_me"];
	        this.can_edit = source["can_edit"];
	        this.is_like = source["is_like"];
	        this.is_open_ledgers = source["is_open_ledgers"];
	        this.is_permission = source["is_permission"];
	        this.is_reposted = source["is_reposted"];
	        this.is_upmost = source["is_upmost"];
	        this.level_permission = source["level_permission"];
	        this.lesson = this.convertValues(source["lesson"], Lesson);
	        this.level = source["level"];
	        this.level_type = source["level_type"];
	        this.log_id = source["log_id"];
	        this.log_type = source["log_type"];
	        this.note = source["note"];
	        this.note_id = source["note_id"];
	        this.note_id_hazy = source["note_id_hazy"];
	        this.note_id_str = source["note_id_str"];
	        this.note_line = source["note_line"];
	        this.note_line_style = source["note_line_style"];
	        this.note_title = source["note_title"];
	        this.note_type = source["note_type"];
	        this.notes_count = this.convertValues(source["notes_count"], NotesCount);
	        this.notes_owner = this.convertValues(source["notes_owner"], NotesUser);
	        this.notes_topic_info = this.convertValues(source["notes_topic_info"], NotesTopicInfo);
	        this.origin_audit_state = source["origin_audit_state"];
	        this.origin_content_type = source["origin_content_type"];
	        this.origin_note_id_hazy = source["origin_note_id_hazy"];
	        this.origin_note_id_str = source["origin_note_id_str"];
	        this.origin_notes_owner = this.convertValues(source["origin_notes_owner"], NotesUser);
	        this.origin_state = source["origin_state"];
	        this.pid = source["pid"];
	        this.pid_str = source["pid_str"];
	        this.ptype = source["ptype"];
	        this.ref_id = source["ref_id"];
	        this.root_audit_state = source["root_audit_state"];
	        this.root_content_type = source["root_content_type"];
	        this.root_highlights = source["root_highlights"];
	        this.root_note_id = source["root_note_id"];
	        this.root_note_id_hazy = source["root_note_id_hazy"];
	        this.root_note_id_str = source["root_note_id_str"];
	        this.root_notes_owner = this.convertValues(source["root_notes_owner"], NotesUser);
	        this.root_state = source["root_state"];
	        this.share_url = source["share_url"];
	        this.source_type = source["source_type"];
	        this.state = source["state"];
	        this.style_note_line = source["style_note_line"];
	        this.switch = this.convertValues(source["switch"], Switch);
	        this.tags = source["tags"];
	        this.uid = source["uid"];
	        this.uid_hazy = source["uid_hazy"];
	        this.update_time = source["update_time"];
	        this.user_expect_status = source["user_expect_status"];
	        this.video = this.convertValues(source["video"], CommentVideo);
	        this.comment_id_str = source["comment_id_str"];
	        this.comment_reply_user = this.convertValues(source["comment_reply_user"], NotesUser);
	        this.repost_comment_note_id_str = source["repost_comment_note_id_str"];
	        this.repost_comment_note_id_hazy = source["repost_comment_note_id_hazy"];
	        this.pid_hazy = source["pid_hazy"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class WendaExtInfo {
	    answer_id: number;
	
	    static createFrom(source: any = {}) {
	        return new WendaExtInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.answer_id = source["answer_id"];
	    }
	}
	export class ReplierInfo {
	    replier_uid: number;
	    replier_name: string;
	    replier_img: string;
	    replier_intro: string;
	    replier_v_status: boolean;
	    replier_v_state_value: number;
	    replier_title: string;
	
	    static createFrom(source: any = {}) {
	        return new ReplierInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.replier_uid = source["replier_uid"];
	        this.replier_name = source["replier_name"];
	        this.replier_img = source["replier_img"];
	        this.replier_intro = source["replier_intro"];
	        this.replier_v_status = source["replier_v_status"];
	        this.replier_v_state_value = source["replier_v_state_value"];
	        this.replier_title = source["replier_title"];
	    }
	}
	export class Course {
	    enid: string;
	    id: number;
	    type: number;
	    class_type: number;
	    class_id: number;
	    has_extra: boolean;
	    class_finished: boolean;
	    title: string;
	    intro: string;
	    author: string;
	    icon: string;
	    create_time: number;
	    last_read: string;
	    progress: number;
	    duration: number;
	    course_num: number;
	    publish_num: number;
	    log_id: string;
	    log_type: string;
	    is_top: number;
	    last_action_time: number;
	    is_new: number;
	    is_finished: number;
	    size: string;
	    dd_url: string;
	    assets_type: number;
	    drm_token: string;
	    audio_detail: Audio;
	    product_price: number;
	    price: string;
	    product_intro: string;
	    has_play_auth: boolean;
	    ext_info: ReplierInfo[];
	    status: number;
	    dd_ext_url: string;
	    is_collected: boolean;
	    wenda_ext_info: WendaExtInfo;
	
	    static createFrom(source: any = {}) {
	        return new Course(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enid = source["enid"];
	        this.id = source["id"];
	        this.type = source["type"];
	        this.class_type = source["class_type"];
	        this.class_id = source["class_id"];
	        this.has_extra = source["has_extra"];
	        this.class_finished = source["class_finished"];
	        this.title = source["title"];
	        this.intro = source["intro"];
	        this.author = source["author"];
	        this.icon = source["icon"];
	        this.create_time = source["create_time"];
	        this.last_read = source["last_read"];
	        this.progress = source["progress"];
	        this.duration = source["duration"];
	        this.course_num = source["course_num"];
	        this.publish_num = source["publish_num"];
	        this.log_id = source["log_id"];
	        this.log_type = source["log_type"];
	        this.is_top = source["is_top"];
	        this.last_action_time = source["last_action_time"];
	        this.is_new = source["is_new"];
	        this.is_finished = source["is_finished"];
	        this.size = source["size"];
	        this.dd_url = source["dd_url"];
	        this.assets_type = source["assets_type"];
	        this.drm_token = source["drm_token"];
	        this.audio_detail = this.convertValues(source["audio_detail"], Audio);
	        this.product_price = source["product_price"];
	        this.price = source["price"];
	        this.product_intro = source["product_intro"];
	        this.has_play_auth = source["has_play_auth"];
	        this.ext_info = this.convertValues(source["ext_info"], ReplierInfo);
	        this.status = source["status"];
	        this.dd_ext_url = source["dd_ext_url"];
	        this.is_collected = source["is_collected"];
	        this.wenda_ext_info = this.convertValues(source["wenda_ext_info"], WendaExtInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CourseCategory {
	    name: string;
	    count: number;
	    category: string;
	
	    static createFrom(source: any = {}) {
	        return new CourseCategory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.count = source["count"];
	        this.category = source["category"];
	    }
	}
	export class LiveInnerArticleInfo {
	    article_id: number;
	    article_title: string;
	    log_type: string;
	    log_id: string;
	
	    static createFrom(source: any = {}) {
	        return new LiveInnerArticleInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.article_id = source["article_id"];
	        this.article_title = source["article_title"];
	        this.log_type = source["log_type"];
	        this.log_id = source["log_id"];
	    }
	}
	export class CourseIntro {
	    type: number;
	    title: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new CourseIntro(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.title = source["title"];
	        this.content = source["content"];
	    }
	}
	export class CourseInfo {
	    class_info: ClassInfo;
	    items: CourseIntro[];
	    intro_article: ArticleIntro;
	    chapter_list: Chapter[];
	    flat_article_list: ArticleBase[];
	    user_type: string;
	    has_more_flat_article_list: boolean;
	    is_show_grading: boolean;
	    class_video: any;
	    live_info: any;
	    time_now: number;
	    class_reviews_count: number;
	    class_reviews: ClassReviews;
	    achievement_detail: any;
	    class_comment_info: ClassCommentInfo;
	    live_inner_article_info: LiveInnerArticleInfo;
	
	    static createFrom(source: any = {}) {
	        return new CourseInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.class_info = this.convertValues(source["class_info"], ClassInfo);
	        this.items = this.convertValues(source["items"], CourseIntro);
	        this.intro_article = this.convertValues(source["intro_article"], ArticleIntro);
	        this.chapter_list = this.convertValues(source["chapter_list"], Chapter);
	        this.flat_article_list = this.convertValues(source["flat_article_list"], ArticleBase);
	        this.user_type = source["user_type"];
	        this.has_more_flat_article_list = source["has_more_flat_article_list"];
	        this.is_show_grading = source["is_show_grading"];
	        this.class_video = source["class_video"];
	        this.live_info = source["live_info"];
	        this.time_now = source["time_now"];
	        this.class_reviews_count = source["class_reviews_count"];
	        this.class_reviews = this.convertValues(source["class_reviews"], ClassReviews);
	        this.achievement_detail = source["achievement_detail"];
	        this.class_comment_info = this.convertValues(source["class_comment_info"], ClassCommentInfo);
	        this.live_inner_article_info = this.convertValues(source["live_inner_article_info"], LiveInnerArticleInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class CourseList {
	    list: Course[];
	    is_more: number;
	    page: number;
	
	    static createFrom(source: any = {}) {
	        return new CourseList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.list = this.convertValues(source["list"], Course);
	        this.is_more = source["is_more"];
	        this.page = source["page"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class SelfInfo {
	    score: number;
	    audit_state: number;
	
	    static createFrom(source: any = {}) {
	        return new SelfInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.score = source["score"];
	        this.audit_state = source["audit_state"];
	    }
	}
	export class EbookScore {
	    id: number;
	    pid: number;
	    ptype: number;
	    isbn: string;
	    average_score: string;
	    score_info: string[];
	    create_time: string;
	    update_time: string;
	    total: string;
	    status: number;
	    book_status: number;
	
	    static createFrom(source: any = {}) {
	        return new EbookScore(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.pid = source["pid"];
	        this.ptype = source["ptype"];
	        this.isbn = source["isbn"];
	        this.average_score = source["average_score"];
	        this.score_info = source["score_info"];
	        this.create_time = source["create_time"];
	        this.update_time = source["update_time"];
	        this.total = source["total"];
	        this.status = source["status"];
	        this.book_status = source["book_status"];
	    }
	}
	export class EbookCommentList {
	    ebook_score: EbookScore;
	    self_info: SelfInfo;
	    list: Comment[];
	    total: number;
	    book: Book;
	
	    static createFrom(source: any = {}) {
	        return new EbookCommentList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ebook_score = this.convertValues(source["ebook_score"], EbookScore);
	        this.self_info = this.convertValues(source["self_info"], SelfInfo);
	        this.list = this.convertValues(source["list"], Comment);
	        this.total = source["total"];
	        this.book = this.convertValues(source["book"], Book);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Press {
	    name: string;
	    brief: string;
	
	    static createFrom(source: any = {}) {
	        return new Press(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.brief = source["brief"];
	    }
	}
	export class EbookDetail {
	    id: number;
	    title: string;
	    style: number;
	    cover: string;
	    count: number;
	    price: string;
	    status: number;
	    operating_title: string;
	    other_share_title: string;
	    other_share_summary: string;
	    author_info: string;
	    book_author: string;
	    publish_time: string;
	    catalog_list: Catalog[];
	    book_intro: string;
	    b_special_price: string;
	    current_price: string;
	    is_buy: boolean;
	    is_trial: boolean;
	    is_tts_switch: boolean;
	    log_id: string;
	    log_type: string;
	    original_price: string;
	    author_list: string[];
	    can_trial_read: boolean;
	    trial_read_proportion: string;
	    with_video: boolean;
	    enid: string;
	    b_overseas_purchase: number;
	    rank_name: string;
	    rank_num: number;
	    is_vip_book: number;
	    is_on_bookshelf: boolean;
	    product_score: string;
	    read_time: number;
	    read_number: any[];
	    press: Press;
	    douban_score: string;
	    classify_name: string;
	    classify_id: number;
	    add_studylist_dd_url: string;
	
	    static createFrom(source: any = {}) {
	        return new EbookDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.style = source["style"];
	        this.cover = source["cover"];
	        this.count = source["count"];
	        this.price = source["price"];
	        this.status = source["status"];
	        this.operating_title = source["operating_title"];
	        this.other_share_title = source["other_share_title"];
	        this.other_share_summary = source["other_share_summary"];
	        this.author_info = source["author_info"];
	        this.book_author = source["book_author"];
	        this.publish_time = source["publish_time"];
	        this.catalog_list = this.convertValues(source["catalog_list"], Catalog);
	        this.book_intro = source["book_intro"];
	        this.b_special_price = source["b_special_price"];
	        this.current_price = source["current_price"];
	        this.is_buy = source["is_buy"];
	        this.is_trial = source["is_trial"];
	        this.is_tts_switch = source["is_tts_switch"];
	        this.log_id = source["log_id"];
	        this.log_type = source["log_type"];
	        this.original_price = source["original_price"];
	        this.author_list = source["author_list"];
	        this.can_trial_read = source["can_trial_read"];
	        this.trial_read_proportion = source["trial_read_proportion"];
	        this.with_video = source["with_video"];
	        this.enid = source["enid"];
	        this.b_overseas_purchase = source["b_overseas_purchase"];
	        this.rank_name = source["rank_name"];
	        this.rank_num = source["rank_num"];
	        this.is_vip_book = source["is_vip_book"];
	        this.is_on_bookshelf = source["is_on_bookshelf"];
	        this.product_score = source["product_score"];
	        this.read_time = source["read_time"];
	        this.read_number = source["read_number"];
	        this.press = this.convertValues(source["press"], Press);
	        this.douban_score = source["douban_score"];
	        this.classify_name = source["classify_name"];
	        this.classify_id = source["classify_id"];
	        this.add_studylist_dd_url = source["add_studylist_dd_url"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class EbookShelfAddData {
	    n: number;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new EbookShelfAddData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.n = source["n"];
	        this.count = source["count"];
	    }
	}
	export class EbookShelfAddResp {
	    data: EbookShelfAddData;
	
	    static createFrom(source: any = {}) {
	        return new EbookShelfAddResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], EbookShelfAddData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class EbookVIPInfo {
	    uid: number;
	    nickname: string;
	    slogan: string;
	    avatar: string;
	    avatar_s: string;
	    month_count: number;
	    total_count: number;
	    week_count: number;
	    finished_count: number;
	    save_price: string;
	    is_vip: boolean;
	    begin_time: number;
	    end_time: number;
	    enterprise_end_time: number;
	    expire_time: number;
	    surplus_time: number;
	    is_expire: boolean;
	    card_id: number;
	    card_type: number;
	    price_desc: string;
	    is_buy_month_discount: boolean;
	    month_discount_price: number;
	    dd_url: string;
	    err_tips: string;
	    v_state_value: number;
	    upgrade_tips: any[];
	
	    static createFrom(source: any = {}) {
	        return new EbookVIPInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uid = source["uid"];
	        this.nickname = source["nickname"];
	        this.slogan = source["slogan"];
	        this.avatar = source["avatar"];
	        this.avatar_s = source["avatar_s"];
	        this.month_count = source["month_count"];
	        this.total_count = source["total_count"];
	        this.week_count = source["week_count"];
	        this.finished_count = source["finished_count"];
	        this.save_price = source["save_price"];
	        this.is_vip = source["is_vip"];
	        this.begin_time = source["begin_time"];
	        this.end_time = source["end_time"];
	        this.enterprise_end_time = source["enterprise_end_time"];
	        this.expire_time = source["expire_time"];
	        this.surplus_time = source["surplus_time"];
	        this.is_expire = source["is_expire"];
	        this.card_id = source["card_id"];
	        this.card_type = source["card_type"];
	        this.price_desc = source["price_desc"];
	        this.is_buy_month_discount = source["is_buy_month_discount"];
	        this.month_discount_price = source["month_discount_price"];
	        this.dd_url = source["dd_url"];
	        this.err_tips = source["err_tips"];
	        this.v_state_value = source["v_state_value"];
	        this.upgrade_tips = source["upgrade_tips"];
	    }
	}
	export class HTab {
	    id: number;
	    title: string;
	    subTitle: string;
	    searchKey: string;
	    url: string;
	    pinnedDeadLine: any;
	    sceneId: number;
	    uv: number;
	    click: number;
	    ctr: number;
	    trackInfo: string;
	    turnType: number;
	    hotType: number;
	    log_id: number;
	    log_type: string;
	
	    static createFrom(source: any = {}) {
	        return new HTab(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.subTitle = source["subTitle"];
	        this.searchKey = source["searchKey"];
	        this.url = source["url"];
	        this.pinnedDeadLine = source["pinnedDeadLine"];
	        this.sceneId = source["sceneId"];
	        this.uv = source["uv"];
	        this.click = source["click"];
	        this.ctr = source["ctr"];
	        this.trackInfo = source["trackInfo"];
	        this.turnType = source["turnType"];
	        this.hotType = source["hotType"];
	        this.log_id = source["log_id"];
	        this.log_type = source["log_type"];
	    }
	}
	export class Label {
	    enid: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Label(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enid = source["enid"];
	        this.name = source["name"];
	    }
	}
	export class HomeCategory {
	    englishName: string;
	    enid: string;
	    icon: string;
	    id: number;
	    labelList: Label[];
	    name: string;
	    navType: number;
	    relationId: number;
	    relationName: string;
	    type: number;
	
	    static createFrom(source: any = {}) {
	        return new HomeCategory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.englishName = source["englishName"];
	        this.enid = source["enid"];
	        this.icon = source["icon"];
	        this.id = source["id"];
	        this.labelList = this.convertValues(source["labelList"], Label);
	        this.name = source["name"];
	        this.navType = source["navType"];
	        this.relationId = source["relationId"];
	        this.relationName = source["relationName"];
	        this.type = source["type"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class HomeModule {
	    description: string;
	    ext1: string;
	    ext2: string;
	    ext3: string;
	    ext4: string;
	    ext5: string;
	    id: number;
	    isShow: number;
	    name: string;
	    sort: number;
	    title: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new HomeModule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.description = source["description"];
	        this.ext1 = source["ext1"];
	        this.ext2 = source["ext2"];
	        this.ext3 = source["ext3"];
	        this.ext4 = source["ext4"];
	        this.ext5 = source["ext5"];
	        this.id = source["id"];
	        this.isShow = source["isShow"];
	        this.name = source["name"];
	        this.sort = source["sort"];
	        this.title = source["title"];
	        this.type = source["type"];
	    }
	}
	export class HomeData {
	    moduleList: HomeModule[];
	    categoryList: HomeCategory[];
	    banner: Banner[];
	
	    static createFrom(source: any = {}) {
	        return new HomeData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.moduleList = this.convertValues(source["moduleList"], HomeModule);
	        this.categoryList = this.convertValues(source["categoryList"], HomeCategory);
	        this.banner = this.convertValues(source["banner"], Banner);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class HomeInitState {
	    isLogin: boolean;
	    homeData: HomeData;
	    uid: string;
	
	    static createFrom(source: any = {}) {
	        return new HomeInitState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.isLogin = source["isLogin"];
	        this.homeData = this.convertValues(source["homeData"], HomeData);
	        this.uid = source["uid"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class HotTab {
	    id: number;
	    name: string;
	    scene_id: number;
	    status: number;
	    sort_num: number;
	    // Go type: time
	    create_time: any;
	    // Go type: time
	    update_time: any;
	    scene_name: string;
	    list: HTab[];
	
	    static createFrom(source: any = {}) {
	        return new HotTab(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.scene_id = source["scene_id"];
	        this.status = source["status"];
	        this.sort_num = source["sort_num"];
	        this.create_time = this.convertValues(source["create_time"], null);
	        this.update_time = this.convertValues(source["update_time"], null);
	        this.scene_name = source["scene_name"];
	        this.list = this.convertValues(source["list"], HTab);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	export class VolcFormat {
	    type: string;
	    format: string;
	    volc_id: string;
	    volc_play_auth_token: string;
	    volc_key_token: string;
	
	    static createFrom(source: any = {}) {
	        return new VolcFormat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.format = source["format"];
	        this.volc_id = source["volc_id"];
	        this.volc_play_auth_token = source["volc_play_auth_token"];
	        this.volc_key_token = source["volc_key_token"];
	    }
	}
	export class VideoTrack {
	    track_id: number;
	    track_type: number;
	    track_type_tag: string;
	    duration: number;
	    formats: VolcFormat[];
	
	    static createFrom(source: any = {}) {
	        return new VideoTrack(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.track_id = source["track_id"];
	        this.track_type = source["track_type"];
	        this.track_type_tag = source["track_type_tag"];
	        this.duration = source["duration"];
	        this.formats = this.convertValues(source["formats"], VolcFormat);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MediaVolc {
	    media_alias_id: string;
	    // Go type: time
	    last_modify: any;
	    version_id: number;
	    tracks: VideoTrack[];
	
	    static createFrom(source: any = {}) {
	        return new MediaVolc(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.media_alias_id = source["media_alias_id"];
	        this.last_modify = this.convertValues(source["last_modify"], null);
	        this.version_id = source["version_id"];
	        this.tracks = this.convertValues(source["tracks"], VideoTrack);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Navigation {
	    enid: string;
	    id: number;
	    relation_id: number;
	    relation_name: string;
	    type: number;
	    nav_type: number;
	    name: string;
	    icon: string;
	    english_name: string;
	    label_list: Label[];
	
	    static createFrom(source: any = {}) {
	        return new Navigation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enid = source["enid"];
	        this.id = source["id"];
	        this.relation_id = source["relation_id"];
	        this.relation_name = source["relation_name"];
	        this.type = source["type"];
	        this.nav_type = source["nav_type"];
	        this.name = source["name"];
	        this.icon = source["icon"];
	        this.english_name = source["english_name"];
	        this.label_list = this.convertValues(source["label_list"], Label);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class NoteBaseSource {
	    title: string;
	    sub_title: string;
	    img: string;
	    p_type: number;
	    pid_str: string;
	    note_ptype: number;
	    is_pop_login_view: boolean;
	    need_check_buy: boolean;
	    url1: string;
	    url2: string;
	    camp_hazy: CampHazy;
	
	    static createFrom(source: any = {}) {
	        return new NoteBaseSource(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.sub_title = source["sub_title"];
	        this.img = source["img"];
	        this.p_type = source["p_type"];
	        this.pid_str = source["pid_str"];
	        this.note_ptype = source["note_ptype"];
	        this.is_pop_login_view = source["is_pop_login_view"];
	        this.need_check_buy = source["need_check_buy"];
	        this.url1 = source["url1"];
	        this.url2 = source["url2"];
	        this.camp_hazy = this.convertValues(source["camp_hazy"], CampHazy);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class NoteCount {
	    repost_count: number;
	    comment_count: number;
	    like_count: number;
	
	    static createFrom(source: any = {}) {
	        return new NoteCount(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.repost_count = source["repost_count"];
	        this.comment_count = source["comment_count"];
	        this.like_count = source["like_count"];
	    }
	}
	export class NoteFPart {
	    uid: number;
	    uid_hazy: string;
	    nick_name: string;
	    avatar: string;
	    follow: number;
	    is_v: number;
	    slogan: string;
	    v_info: string;
	    student_id: number;
	    student_id_hazy: string;
	    is_poster: boolean;
	    qr_code: string;
	    note: string;
	    time_desc: string;
	    note_title: string;
	    note_score: string;
	    note_line: string;
	    note_id: string;
	    note_id_hazy: string;
	    tip: string;
	    images: string[];
	    base_source: NoteBaseSource;
	    style_note_line: string;
	
	    static createFrom(source: any = {}) {
	        return new NoteFPart(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uid = source["uid"];
	        this.uid_hazy = source["uid_hazy"];
	        this.nick_name = source["nick_name"];
	        this.avatar = source["avatar"];
	        this.follow = source["follow"];
	        this.is_v = source["is_v"];
	        this.slogan = source["slogan"];
	        this.v_info = source["v_info"];
	        this.student_id = source["student_id"];
	        this.student_id_hazy = source["student_id_hazy"];
	        this.is_poster = source["is_poster"];
	        this.qr_code = source["qr_code"];
	        this.note = source["note"];
	        this.time_desc = source["time_desc"];
	        this.note_title = source["note_title"];
	        this.note_score = source["note_score"];
	        this.note_line = source["note_line"];
	        this.note_id = source["note_id"];
	        this.note_id_hazy = source["note_id_hazy"];
	        this.tip = source["tip"];
	        this.images = source["images"];
	        this.base_source = this.convertValues(source["base_source"], NoteBaseSource);
	        this.style_note_line = source["style_note_line"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class NoteTopic {
	    topic_id: number;
	    topic_id_hazy: string;
	    is_elected: boolean;
	    is_topmost: boolean;
	    topic_name: string;
	
	    static createFrom(source: any = {}) {
	        return new NoteTopic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.topic_id = source["topic_id"];
	        this.topic_id_hazy = source["topic_id_hazy"];
	        this.is_elected = source["is_elected"];
	        this.is_topmost = source["is_topmost"];
	        this.topic_name = source["topic_name"];
	    }
	}
	export class NoteDetail {
	    detail_title: string;
	    comb?: Comb[];
	    state: number;
	    is_mine: boolean;
	    is_reposted: boolean;
	    is_like: boolean;
	    own_uid_hazy: string;
	    topic: NoteTopic;
	    tags: any[];
	    folders: any;
	    note_count: NoteCount;
	    f_part: NoteFPart;
	    s_part?: NoteFPart;
	    share_url: string;
	    class: number;
	    level: number;
	    level_type: number;
	    level_permission: boolean;
	    notes_type: number;
	
	    static createFrom(source: any = {}) {
	        return new NoteDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.detail_title = source["detail_title"];
	        this.comb = this.convertValues(source["comb"], Comb);
	        this.state = source["state"];
	        this.is_mine = source["is_mine"];
	        this.is_reposted = source["is_reposted"];
	        this.is_like = source["is_like"];
	        this.own_uid_hazy = source["own_uid_hazy"];
	        this.topic = this.convertValues(source["topic"], NoteTopic);
	        this.tags = source["tags"];
	        this.folders = source["folders"];
	        this.note_count = this.convertValues(source["note_count"], NoteCount);
	        this.f_part = this.convertValues(source["f_part"], NoteFPart);
	        this.s_part = this.convertValues(source["s_part"], NoteFPart);
	        this.share_url = source["share_url"];
	        this.class = source["class"];
	        this.level = source["level"];
	        this.level_type = source["level_type"];
	        this.level_permission = source["level_permission"];
	        this.notes_type = source["notes_type"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	export class NotesList {
	    has_more: boolean;
	    list: any[];
	    note_detail_list: NoteDetail[];
	    presenter_uid: any[];
	    presenter_uid_hazy: string[];
	
	    static createFrom(source: any = {}) {
	        return new NotesList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.has_more = source["has_more"];
	        this.list = source["list"];
	        this.note_detail_list = this.convertValues(source["note_detail_list"], NoteDetail);
	        this.presenter_uid = source["presenter_uid"];
	        this.presenter_uid_hazy = source["presenter_uid_hazy"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class NotesTimeline {
	    is_more: boolean;
	    max_id: string;
	    notes: NoteDetail[];
	
	    static createFrom(source: any = {}) {
	        return new NotesTimeline(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.is_more = source["is_more"];
	        this.max_id = source["max_id"];
	        this.notes = this.convertValues(source["notes"], NoteDetail);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class OdobRight {
	    name: string;
	    right: boolean;
	
	    static createFrom(source: any = {}) {
	        return new OdobRight(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.right = source["right"];
	    }
	}
	export class OdobCard {
	    available_buy: boolean;
	    can_gift: number;
	    card_type: number;
	    contract_template: string;
	    cover: string;
	    created_at: string;
	    days: number;
	    description: string;
	    device: string;
	    discount_tip: string;
	    displayed: number;
	    first_month_price: string;
	    id: number;
	    ios_switch: boolean;
	    is_first_subscribe: boolean;
	    is_subscribed: number;
	    max_num: number;
	    name: string;
	    origin_price: string;
	    price: string;
	    price_desc: string;
	    ptype: number;
	    rights: OdobRight[];
	    selected: number;
	    sku: string;
	    sort_no: number;
	    status: number;
	    style_type: number;
	    subscribe_desc: string;
	    subscribed_flag: string;
	    type: number;
	    updated_at: string;
	    used_card_count: number;
	    user_max_buy_count: number;
	    welfare_info: string;
	
	    static createFrom(source: any = {}) {
	        return new OdobCard(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.available_buy = source["available_buy"];
	        this.can_gift = source["can_gift"];
	        this.card_type = source["card_type"];
	        this.contract_template = source["contract_template"];
	        this.cover = source["cover"];
	        this.created_at = source["created_at"];
	        this.days = source["days"];
	        this.description = source["description"];
	        this.device = source["device"];
	        this.discount_tip = source["discount_tip"];
	        this.displayed = source["displayed"];
	        this.first_month_price = source["first_month_price"];
	        this.id = source["id"];
	        this.ios_switch = source["ios_switch"];
	        this.is_first_subscribe = source["is_first_subscribe"];
	        this.is_subscribed = source["is_subscribed"];
	        this.max_num = source["max_num"];
	        this.name = source["name"];
	        this.origin_price = source["origin_price"];
	        this.price = source["price"];
	        this.price_desc = source["price_desc"];
	        this.ptype = source["ptype"];
	        this.rights = this.convertValues(source["rights"], OdobRight);
	        this.selected = source["selected"];
	        this.sku = source["sku"];
	        this.sort_no = source["sort_no"];
	        this.status = source["status"];
	        this.style_type = source["style_type"];
	        this.subscribe_desc = source["subscribe_desc"];
	        this.subscribed_flag = source["subscribed_flag"];
	        this.type = source["type"];
	        this.updated_at = source["updated_at"];
	        this.used_card_count = source["used_card_count"];
	        this.user_max_buy_count = source["user_max_buy_count"];
	        this.welfare_info = source["welfare_info"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class OdobUser {
	    avatar: string;
	    avatar_s: string;
	    begin_time: number;
	    can_given_count: number;
	    card_id: number;
	    dd_url: string;
	    end_time: number;
	    enterprise_end_time: number;
	    err_tips: string;
	    expire_time: number;
	    given_card_count: number;
	    is_expire: boolean;
	    is_vip: boolean;
	    nickname: string;
	    price_desc: string;
	    save_price: string;
	    slogan: string;
	    surplus_time: number;
	    total_count: number;
	    uid: number;
	    v_state_value: number;
	    week_count: number;
	
	    static createFrom(source: any = {}) {
	        return new OdobUser(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.avatar = source["avatar"];
	        this.avatar_s = source["avatar_s"];
	        this.begin_time = source["begin_time"];
	        this.can_given_count = source["can_given_count"];
	        this.card_id = source["card_id"];
	        this.dd_url = source["dd_url"];
	        this.end_time = source["end_time"];
	        this.enterprise_end_time = source["enterprise_end_time"];
	        this.err_tips = source["err_tips"];
	        this.expire_time = source["expire_time"];
	        this.given_card_count = source["given_card_count"];
	        this.is_expire = source["is_expire"];
	        this.is_vip = source["is_vip"];
	        this.nickname = source["nickname"];
	        this.price_desc = source["price_desc"];
	        this.save_price = source["save_price"];
	        this.slogan = source["slogan"];
	        this.surplus_time = source["surplus_time"];
	        this.total_count = source["total_count"];
	        this.uid = source["uid"];
	        this.v_state_value = source["v_state_value"];
	        this.week_count = source["week_count"];
	    }
	}
	export class OdobVip {
	    card: OdobCard[];
	    user: OdobUser;
	
	    static createFrom(source: any = {}) {
	        return new OdobVip(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.card = this.convertValues(source["card"], OdobCard);
	        this.user = this.convertValues(source["user"], OdobUser);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class ProductSimple {
	    product_type: number;
	    product_enid: string;
	    title: string;
	    intro: string;
	    introduction: string;
	    index_image: string;
	    score: string;
	    user_score_count: number;
	    horizontal_image: string;
	    learn_user_count: number;
	    author_list: string[];
	    trackinfo: string;
	    log_type: string;
	
	    static createFrom(source: any = {}) {
	        return new ProductSimple(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.product_type = source["product_type"];
	        this.product_enid = source["product_enid"];
	        this.title = source["title"];
	        this.intro = source["intro"];
	        this.introduction = source["introduction"];
	        this.index_image = source["index_image"];
	        this.score = source["score"];
	        this.user_score_count = source["user_score_count"];
	        this.horizontal_image = source["horizontal_image"];
	        this.learn_user_count = source["learn_user_count"];
	        this.author_list = source["author_list"];
	        this.trackinfo = source["trackinfo"];
	        this.log_type = source["log_type"];
	    }
	}
	
	
	export class RTab {
	    id: number;
	    idstr: string;
	    name: string;
	    rank: number;
	    status: number;
	    create_time: number;
	    update_time: number;
	    finish_time: number;
	    type: number;
	    is_hot: number;
	    is_top: number;
	    track_info: string;
	    search_key: string;
	
	    static createFrom(source: any = {}) {
	        return new RTab(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.idstr = source["idstr"];
	        this.name = source["name"];
	        this.rank = source["rank"];
	        this.status = source["status"];
	        this.create_time = source["create_time"];
	        this.update_time = source["update_time"];
	        this.finish_time = source["finish_time"];
	        this.type = source["type"];
	        this.is_hot = source["is_hot"];
	        this.is_top = source["is_top"];
	        this.track_info = source["track_info"];
	        this.search_key = source["search_key"];
	    }
	}
	
	export class Recommend {
	    ptype_scene_name: string;
	    type: number;
	    list: RTab[];
	
	    static createFrom(source: any = {}) {
	        return new Recommend(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ptype_scene_name = source["ptype_scene_name"];
	        this.type = source["type"];
	        this.list = this.convertValues(source["list"], RTab);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class SearchTot {
	    hot_tab_list: HotTab[];
	    recommend_map: Recommend[];
	
	    static createFrom(source: any = {}) {
	        return new SearchTot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hot_tab_list = this.convertValues(source["hot_tab_list"], HotTab);
	        this.recommend_map = this.convertValues(source["recommend_map"], Recommend);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class SunflowerContent {
	    product_list: ProductSimple[];
	    current_enid: string;
	    navigation_list: Navigation[];
	    page_id: number;
	    page_size: number;
	    is_more: number;
	    request_id: string;
	
	    static createFrom(source: any = {}) {
	        return new SunflowerContent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.product_list = this.convertValues(source["product_list"], ProductSimple);
	        this.current_enid = source["current_enid"];
	        this.navigation_list = this.convertValues(source["navigation_list"], Navigation);
	        this.page_id = source["page_id"];
	        this.page_size = source["page_size"];
	        this.is_more = source["is_more"];
	        this.request_id = source["request_id"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SunflowerLabelList {
	    list: Navigation[];
	
	    static createFrom(source: any = {}) {
	        return new SunflowerLabelList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.list = this.convertValues(source["list"], Navigation);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SunflowerResource {
	    id: number;
	    enid: string;
	    name: string;
	    intro: string;
	    logo: string;
	    product_type: number;
	    product_id: number;
	    score: number;
	    class_type: number;
	    status: number;
	
	    static createFrom(source: any = {}) {
	        return new SunflowerResource(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.enid = source["enid"];
	        this.name = source["name"];
	        this.intro = source["intro"];
	        this.logo = source["logo"];
	        this.product_type = source["product_type"];
	        this.product_id = source["product_id"];
	        this.score = source["score"];
	        this.class_type = source["class_type"];
	        this.status = source["status"];
	    }
	}
	export class SunflowerResourceList {
	    list: SunflowerResource[];
	
	    static createFrom(source: any = {}) {
	        return new SunflowerResourceList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.list = this.convertValues(source["list"], SunflowerResource);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class TopicIntro {
	    notes_topic_id: string;
	    topic_id_hazy: string;
	    name: string;
	    img: string;
	    topmost: boolean;
	    tag: number;
	    intro: string;
	    view_count: number;
	    notes_count: number;
	    has_new_notes: boolean;
	    user_state: number;
	    log_id: string;
	    log_type: string;
	
	    static createFrom(source: any = {}) {
	        return new TopicIntro(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.notes_topic_id = source["notes_topic_id"];
	        this.topic_id_hazy = source["topic_id_hazy"];
	        this.name = source["name"];
	        this.img = source["img"];
	        this.topmost = source["topmost"];
	        this.tag = source["tag"];
	        this.intro = source["intro"];
	        this.view_count = source["view_count"];
	        this.notes_count = source["notes_count"];
	        this.has_new_notes = source["has_new_notes"];
	        this.user_state = source["user_state"];
	        this.log_id = source["log_id"];
	        this.log_type = source["log_type"];
	    }
	}
	export class TopicAll {
	    has_more: boolean;
	    list: TopicIntro[];
	
	    static createFrom(source: any = {}) {
	        return new TopicAll(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.has_more = source["has_more"];
	        this.list = this.convertValues(source["list"], TopicIntro);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TopicTopArea {
	    id: number;
	    id_hazy: string;
	    icon: string;
	    title: string;
	    url: string;
	    state: number;
	    index_num: number;
	
	    static createFrom(source: any = {}) {
	        return new TopicTopArea(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.id_hazy = source["id_hazy"];
	        this.icon = source["icon"];
	        this.title = source["title"];
	        this.url = source["url"];
	        this.state = source["state"];
	        this.index_num = source["index_num"];
	    }
	}
	export class TopicPresenter {
	    id: number;
	    id_hazy: string;
	    uid: number;
	    uid_hazy: string;
	    isV: number;
	    name: string;
	    avatar: string;
	    topic_count: number;
	    icon: string;
	    relation: number;
	
	    static createFrom(source: any = {}) {
	        return new TopicPresenter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.id_hazy = source["id_hazy"];
	        this.uid = source["uid"];
	        this.uid_hazy = source["uid_hazy"];
	        this.isV = source["isV"];
	        this.name = source["name"];
	        this.avatar = source["avatar"];
	        this.topic_count = source["topic_count"];
	        this.icon = source["icon"];
	        this.relation = source["relation"];
	    }
	}
	export class TopicDetail {
	    notes_topic_id: string;
	    topic_id_hazy: string;
	    name: string;
	    img: string;
	    topmost: boolean;
	    tag: number;
	    intro: string;
	    view_count: number;
	    notes_count: number;
	    has_new_notes: boolean;
	    user_state: number;
	    log_id: string;
	    log_type: string;
	    presenters: TopicPresenter[];
	    top_area: TopicTopArea[];
	    create_time: number;
	    update_time: number;
	    state: number;
	    share_url: string;
	    ddurl: string;
	    topic_id_str: string;
	    last_notes_id: string;
	    last_notes_id_hazy: string;
	    last_notes_uid: string;
	    last_notes_uid_hazy: string;
	    last_notes_content: string;
	    last_update_time: number;
	
	    static createFrom(source: any = {}) {
	        return new TopicDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.notes_topic_id = source["notes_topic_id"];
	        this.topic_id_hazy = source["topic_id_hazy"];
	        this.name = source["name"];
	        this.img = source["img"];
	        this.topmost = source["topmost"];
	        this.tag = source["tag"];
	        this.intro = source["intro"];
	        this.view_count = source["view_count"];
	        this.notes_count = source["notes_count"];
	        this.has_new_notes = source["has_new_notes"];
	        this.user_state = source["user_state"];
	        this.log_id = source["log_id"];
	        this.log_type = source["log_type"];
	        this.presenters = this.convertValues(source["presenters"], TopicPresenter);
	        this.top_area = this.convertValues(source["top_area"], TopicTopArea);
	        this.create_time = source["create_time"];
	        this.update_time = source["update_time"];
	        this.state = source["state"];
	        this.share_url = source["share_url"];
	        this.ddurl = source["ddurl"];
	        this.topic_id_str = source["topic_id_str"];
	        this.last_notes_id = source["last_notes_id"];
	        this.last_notes_id_hazy = source["last_notes_id_hazy"];
	        this.last_notes_uid = source["last_notes_uid"];
	        this.last_notes_uid_hazy = source["last_notes_uid_hazy"];
	        this.last_notes_content = source["last_notes_content"];
	        this.last_update_time = source["last_update_time"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	export class VIPUser {
	    info: string;
	    stat: number;
	
	    static createFrom(source: any = {}) {
	        return new VIPUser(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.info = source["info"];
	        this.stat = source["stat"];
	    }
	}
	export class User {
	    nickname: string;
	    avatar: string;
	    today_study_time: number;
	    study_serial_days: number;
	    IsV: number;
	    vip_user: VIPUser;
	    is_teacher: number;
	    uid_hazy: string;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nickname = source["nickname"];
	        this.avatar = source["avatar"];
	        this.today_study_time = source["today_study_time"];
	        this.study_serial_days = source["study_serial_days"];
	        this.IsV = source["IsV"];
	        this.vip_user = this.convertValues(source["vip_user"], VIPUser);
	        this.is_teacher = source["is_teacher"];
	        this.uid_hazy = source["uid_hazy"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	export class VodAdaptiveInfo {
	    MainPlayUrl: string;
	    BackupPlayUrl: string;
	    AdaptiveType: string;
	
	    static createFrom(source: any = {}) {
	        return new VodAdaptiveInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.MainPlayUrl = source["MainPlayUrl"];
	        this.BackupPlayUrl = source["BackupPlayUrl"];
	        this.AdaptiveType = source["AdaptiveType"];
	    }
	}
	export class VodPlayInfo {
	    FileId: string;
	    Md5: string;
	    FileType: string;
	    Format: string;
	    Codec: string;
	    Definition: string;
	    MainPlayUrl: string;
	    BackupPlayUrl: string;
	    Bitrate: number;
	    Width: number;
	    Height: number;
	    Size: number;
	    CheckInfo: string;
	    IndexRange: string;
	    InitRange: string;
	    PlayAuth: string;
	    PlayAuthId: string;
	    LogoType: string;
	    Quality: string;
	    BarrageMaskOffset: string;
	    Duration: number;
	    KeyFrameAlignment: string;
	    Volume: any;
	
	    static createFrom(source: any = {}) {
	        return new VodPlayInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.FileId = source["FileId"];
	        this.Md5 = source["Md5"];
	        this.FileType = source["FileType"];
	        this.Format = source["Format"];
	        this.Codec = source["Codec"];
	        this.Definition = source["Definition"];
	        this.MainPlayUrl = source["MainPlayUrl"];
	        this.BackupPlayUrl = source["BackupPlayUrl"];
	        this.Bitrate = source["Bitrate"];
	        this.Width = source["Width"];
	        this.Height = source["Height"];
	        this.Size = source["Size"];
	        this.CheckInfo = source["CheckInfo"];
	        this.IndexRange = source["IndexRange"];
	        this.InitRange = source["InitRange"];
	        this.PlayAuth = source["PlayAuth"];
	        this.PlayAuthId = source["PlayAuthId"];
	        this.LogoType = source["LogoType"];
	        this.Quality = source["Quality"];
	        this.BarrageMaskOffset = source["BarrageMaskOffset"];
	        this.Duration = source["Duration"];
	        this.KeyFrameAlignment = source["KeyFrameAlignment"];
	        this.Volume = source["Volume"];
	    }
	}
	export class VodThumbInfo {
	    CaptureNum: number;
	    StoreUrls: string[];
	    CellWidth: number;
	    CellHeight: number;
	    ImgXLen: number;
	    ImgYLen: number;
	    Interval: number;
	    Format: string;
	
	    static createFrom(source: any = {}) {
	        return new VodThumbInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.CaptureNum = source["CaptureNum"];
	        this.StoreUrls = source["StoreUrls"];
	        this.CellWidth = source["CellWidth"];
	        this.CellHeight = source["CellHeight"];
	        this.ImgXLen = source["ImgXLen"];
	        this.ImgYLen = source["ImgYLen"];
	        this.Interval = source["Interval"];
	        this.Format = source["Format"];
	    }
	}
	export class VodPlayInfoModel {
	    Version: number;
	    Vid: string;
	    Status: number;
	    PosterUrl: string;
	    Duration: number;
	    FileType: string;
	    EnableAdaptive: boolean;
	    TotalCount: number;
	    AdaptiveInfo: VodAdaptiveInfo;
	    PlayInfoList: VodPlayInfo[];
	    ThumbInfoList: VodThumbInfo[];
	    BarrageMaskUrl: string;
	    SubtitleInfoList: any[];
	
	    static createFrom(source: any = {}) {
	        return new VodPlayInfoModel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Version = source["Version"];
	        this.Vid = source["Vid"];
	        this.Status = source["Status"];
	        this.PosterUrl = source["PosterUrl"];
	        this.Duration = source["Duration"];
	        this.FileType = source["FileType"];
	        this.EnableAdaptive = source["EnableAdaptive"];
	        this.TotalCount = source["TotalCount"];
	        this.AdaptiveInfo = this.convertValues(source["AdaptiveInfo"], VodAdaptiveInfo);
	        this.PlayInfoList = this.convertValues(source["PlayInfoList"], VodPlayInfo);
	        this.ThumbInfoList = this.convertValues(source["ThumbInfoList"], VodThumbInfo);
	        this.BarrageMaskUrl = source["BarrageMaskUrl"];
	        this.SubtitleInfoList = source["SubtitleInfoList"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class VodRespMetadata {
	    RequestId: string;
	    Action: string;
	    Version: string;
	    Service: string;
	    Region: string;
	
	    static createFrom(source: any = {}) {
	        return new VodRespMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.RequestId = source["RequestId"];
	        this.Action = source["Action"];
	        this.Version = source["Version"];
	        this.Service = source["Service"];
	        this.Region = source["Region"];
	    }
	}
	export class VodPlayInfoResp {
	    ResponseMetadata: VodRespMetadata;
	    Result: VodPlayInfoModel;
	
	    static createFrom(source: any = {}) {
	        return new VodPlayInfoResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ResponseMetadata = this.convertValues(source["ResponseMetadata"], VodRespMetadata);
	        this.Result = this.convertValues(source["Result"], VodPlayInfoModel);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	

}

