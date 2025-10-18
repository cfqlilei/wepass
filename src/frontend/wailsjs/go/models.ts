export namespace models {
	
	export class Account {
	    id: string;
	    title: string;
	    username: string;
	    password: string;
	    url: string;
	    typeid: string;
	    notes: string;
	    icon: string;
	    is_favorite: boolean;
	    use_count: number;
	    // Go type: time
	    last_used_at: any;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    input_method: number;
	
	    static createFrom(source: any = {}) {
	        return new Account(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.url = source["url"];
	        this.typeid = source["typeid"];
	        this.notes = source["notes"];
	        this.icon = source["icon"];
	        this.is_favorite = source["is_favorite"];
	        this.use_count = source["use_count"];
	        this.last_used_at = this.convertValues(source["last_used_at"], null);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.input_method = source["input_method"];
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
	export class AccountDecrypted {
	    id: string;
	    title: string;
	    username: string;
	    password: string;
	    url: string;
	    typeid: string;
	    group_id: string;
	    notes: string;
	    icon: string;
	    is_favorite: boolean;
	    use_count: number;
	    // Go type: time
	    last_used_at: any;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    input_method: number;
	    masked_username: string;
	    masked_password: string;
	
	    static createFrom(source: any = {}) {
	        return new AccountDecrypted(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.url = source["url"];
	        this.typeid = source["typeid"];
	        this.group_id = source["group_id"];
	        this.notes = source["notes"];
	        this.icon = source["icon"];
	        this.is_favorite = source["is_favorite"];
	        this.use_count = source["use_count"];
	        this.last_used_at = this.convertValues(source["last_used_at"], null);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.input_method = source["input_method"];
	        this.masked_username = source["masked_username"];
	        this.masked_password = source["masked_password"];
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
	export class HotkeyConfig {
	    enable_global_hotkey: boolean;
	    show_hide_hotkey: string;
	    enable_show_hide_hotkey: boolean;
	
	    static createFrom(source: any = {}) {
	        return new HotkeyConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enable_global_hotkey = source["enable_global_hotkey"];
	        this.show_hide_hotkey = source["show_hide_hotkey"];
	        this.enable_show_hide_hotkey = source["enable_show_hide_hotkey"];
	    }
	}
	export class LockConfig {
	    enable_auto_lock: boolean;
	    enable_timer_lock: boolean;
	    enable_minimize_lock: boolean;
	    lock_time_minutes: number;
	    enable_system_lock: boolean;
	    system_lock_minutes: number;
	
	    static createFrom(source: any = {}) {
	        return new LockConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enable_auto_lock = source["enable_auto_lock"];
	        this.enable_timer_lock = source["enable_timer_lock"];
	        this.enable_minimize_lock = source["enable_minimize_lock"];
	        this.lock_time_minutes = source["lock_time_minutes"];
	        this.enable_system_lock = source["enable_system_lock"];
	        this.system_lock_minutes = source["system_lock_minutes"];
	    }
	}
	export class LogConfig {
	    enable_info_log: boolean;
	    enable_debug_log: boolean;
	
	    static createFrom(source: any = {}) {
	        return new LogConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enable_info_log = source["enable_info_log"];
	        this.enable_debug_log = source["enable_debug_log"];
	    }
	}
	export class AppConfig {
	    current_vault_path: string;
	    recent_vaults: string[];
	    window_width: number;
	    window_height: number;
	    theme: string;
	    language: string;
	    log_config: LogConfig;
	    lock_config: LockConfig;
	    hotkey_config: HotkeyConfig;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.current_vault_path = source["current_vault_path"];
	        this.recent_vaults = source["recent_vaults"];
	        this.window_width = source["window_width"];
	        this.window_height = source["window_height"];
	        this.theme = source["theme"];
	        this.language = source["language"];
	        this.log_config = this.convertValues(source["log_config"], LogConfig);
	        this.lock_config = this.convertValues(source["lock_config"], LockConfig);
	        this.hotkey_config = this.convertValues(source["hotkey_config"], HotkeyConfig);
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
	export class CustomRuleConfig {
	    pattern: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new CustomRuleConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pattern = source["pattern"];
	        this.description = source["description"];
	    }
	}
	export class GeneralRuleConfig {
	    include_uppercase: boolean;
	    include_lowercase: boolean;
	    include_numbers: boolean;
	    include_special_chars: boolean;
	    include_custom_chars: boolean;
	    min_uppercase: number;
	    min_lowercase: number;
	    min_numbers: number;
	    min_special_chars: number;
	    min_custom_chars: number;
	    length: number;
	    custom_special_chars: string;
	
	    static createFrom(source: any = {}) {
	        return new GeneralRuleConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.include_uppercase = source["include_uppercase"];
	        this.include_lowercase = source["include_lowercase"];
	        this.include_numbers = source["include_numbers"];
	        this.include_special_chars = source["include_special_chars"];
	        this.include_custom_chars = source["include_custom_chars"];
	        this.min_uppercase = source["min_uppercase"];
	        this.min_lowercase = source["min_lowercase"];
	        this.min_numbers = source["min_numbers"];
	        this.min_special_chars = source["min_special_chars"];
	        this.min_custom_chars = source["min_custom_chars"];
	        this.length = source["length"];
	        this.custom_special_chars = source["custom_special_chars"];
	    }
	}
	export class Group {
	    id: string;
	    name: string;
	    icon: string;
	    sort_order: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Group(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.icon = source["icon"];
	        this.sort_order = source["sort_order"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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
	
	
	
	export class PasswordRule {
	    id: string;
	    name: string;
	    description: string;
	    rule_type: string;
	    config: string;
	    is_default: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new PasswordRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.rule_type = source["rule_type"];
	        this.config = source["config"];
	        this.is_default = source["is_default"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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
	export class Type {
	    id: string;
	    name: string;
	    icon: string;
	    filter: string;
	    group_id: string;
	    sort_order: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Type(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.icon = source["icon"];
	        this.filter = source["filter"];
	        this.group_id = source["group_id"];
	        this.sort_order = source["sort_order"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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

export namespace services {
	
	export class SkippedAccountInfo {
	    id: string;
	    title: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new SkippedAccountInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.name = source["name"];
	    }
	}
	export class ImportResult {
	    success: boolean;
	    total_accounts: number;
	    imported_accounts: number;
	    skipped_accounts: number;
	    error_accounts: number;
	    total_groups: number;
	    imported_groups: number;
	    skipped_groups: number;
	    total_types: number;
	    imported_types: number;
	    skipped_types: number;
	    skipped_account_details: SkippedAccountInfo[];
	    error_message: string;
	
	    static createFrom(source: any = {}) {
	        return new ImportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.total_accounts = source["total_accounts"];
	        this.imported_accounts = source["imported_accounts"];
	        this.skipped_accounts = source["skipped_accounts"];
	        this.error_accounts = source["error_accounts"];
	        this.total_groups = source["total_groups"];
	        this.imported_groups = source["imported_groups"];
	        this.skipped_groups = source["skipped_groups"];
	        this.total_types = source["total_types"];
	        this.imported_types = source["imported_types"];
	        this.skipped_types = source["skipped_types"];
	        this.skipped_account_details = this.convertValues(source["skipped_account_details"], SkippedAccountInfo);
	        this.error_message = source["error_message"];
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

export namespace version {
	
	export class ChangeLogEntry {
	    Version: string;
	    Date: string;
	    Changes: string[];
	
	    static createFrom(source: any = {}) {
	        return new ChangeLogEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Version = source["Version"];
	        this.Date = source["Date"];
	        this.Changes = source["Changes"];
	    }
	}

}

