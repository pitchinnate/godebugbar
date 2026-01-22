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
export type ErrorListProps = typeof __propDef.props;
export type ErrorListEvents = typeof __propDef.events;
export type ErrorListSlots = typeof __propDef.slots;
export default class ErrorList extends SvelteComponent<ErrorListProps, ErrorListEvents, ErrorListSlots> {
}
export {};
