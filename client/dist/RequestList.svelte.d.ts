import { SvelteComponent } from "svelte";
import type { RequestInfo } from './types.js';
declare const __propDef: {
    props: {
        requests: RequestInfo[];
        selectedId: string | undefined;
        onSelect: (request: RequestInfo) => void;
        getMethodColor: (method: string) => string;
    };
    events: {
        [evt: string]: CustomEvent<any>;
    };
    slots: {};
    exports?: {} | undefined;
    bindings?: string | undefined;
};
export type RequestListProps = typeof __propDef.props;
export type RequestListEvents = typeof __propDef.events;
export type RequestListSlots = typeof __propDef.slots;
export default class RequestList extends SvelteComponent<RequestListProps, RequestListEvents, RequestListSlots> {
}
export {};
