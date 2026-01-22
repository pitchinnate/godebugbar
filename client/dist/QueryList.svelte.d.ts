import { SvelteComponent } from "svelte";
import type { RequestInfo } from './types.js';
declare const __propDef: {
    props: {
        requests: RequestInfo[];
    };
    events: {
        [evt: string]: CustomEvent<any>;
    };
    slots: {};
    exports?: {} | undefined;
    bindings?: string | undefined;
};
export type QueryListProps = typeof __propDef.props;
export type QueryListEvents = typeof __propDef.events;
export type QueryListSlots = typeof __propDef.slots;
export default class QueryList extends SvelteComponent<QueryListProps, QueryListEvents, QueryListSlots> {
}
export {};
