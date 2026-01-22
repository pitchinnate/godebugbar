import { SvelteComponent } from "svelte";
import type { RequestInfo } from './types.js';
declare const __propDef: {
    props: {
        request: RequestInfo;
        getMethodColor: (method: string) => string;
    };
    events: {
        [evt: string]: CustomEvent<any>;
    };
    slots: {};
    exports?: {} | undefined;
    bindings?: string | undefined;
};
export type RequestDetailProps = typeof __propDef.props;
export type RequestDetailEvents = typeof __propDef.events;
export type RequestDetailSlots = typeof __propDef.slots;
export default class RequestDetail extends SvelteComponent<RequestDetailProps, RequestDetailEvents, RequestDetailSlots> {
}
export {};
