export namespace logger {
	
	export class Log_Data {
	    fields?: {[key: string]: string};
	
	    static createFrom(source: any = {}) {
	        return new Log_Data(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fields = source["fields"];
	    }
	}
	export class Log {
	    id?: number;
	    message?: string;
	    level?: number;
	    severity?: number;
	    timestamp?: number;
	    origin?: string;
	    source?: string;
	    type?: string;
	    group?: string;
	    tags?: string;
	    data?: Log_Data;
	
	    static createFrom(source: any = {}) {
	        return new Log(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.message = source["message"];
	        this.level = source["level"];
	        this.severity = source["severity"];
	        this.timestamp = source["timestamp"];
	        this.origin = source["origin"];
	        this.source = source["source"];
	        this.type = source["type"];
	        this.group = source["group"];
	        this.tags = source["tags"];
	        this.data = this.convertValues(source["data"], Log_Data);
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

