import { SvelteComponent } from "svelte";
import type { DebugBarConfig } from './types.js';
declare const __propDef: {
    props: {
        config: DebugBarConfig;
    };
    events: {
        [evt: string]: CustomEvent<any>;
    };
    slots: {};
    exports?: {} | undefined;
    bindings?: string | undefined;
};
export type DebugBarProps = typeof __propDef.props;
export type DebugBarEvents = typeof __propDef.events;
export type DebugBarSlots = typeof __propDef.slots;
export default class DebugBar extends SvelteComponent<DebugBarProps, DebugBarEvents, DebugBarSlots> {
}
export {};
