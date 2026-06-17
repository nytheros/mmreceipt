import React from 'react';
import ReactDOM from 'react-dom';
type PluginRegistry = {registerWebSocketEventHandler(event: string, handler: (msg: {data?: ReceiptRecord}) => void): void};
import {ReceiptIndicator} from './components/ReceiptIndicator';
import {receiptStore, ReceiptRecord} from './store/receipt_store';
import {startReadReceiptObserver} from './hooks/useReadReceiptObserver';

function injectCSS(): void {
    if (document.getElementById('rr-styles')) return;
    const style = document.createElement('style');
    style.id = 'rr-styles';
    style.textContent = `.rr-delivered { color: #808080 !important; } .rr-read { color: #2196F3 !important; } .rr-sent { color: inherit; }`;
    document.head.appendChild(style);
}

const roots = new Map<Element, any>();
let stopObserver: (() => void) | undefined;

function postIdFromTime(el: Element): string | undefined { return el.closest('[id^="post_"]')?.id.replace(/^post_/, ''); }

let pendingPosts = new Set<string>();
let batchTimer: number | undefined;

function flushBatch(): void {
    batchTimer = undefined;
    const ids = [...pendingPosts];
    if (ids.length === 0) return;
    pendingPosts = new Set<string>();
    void fetch('/plugins/readreceipt/api/v1/status/batch', {
        method: 'POST', credentials: 'same-origin',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({post_ids: ids}),
    }).then((r) => r.ok ? r.json() : undefined).then((data) => {
        if (!data?.statuses) return;
        for (const [postId, receipts] of Object.entries(data.statuses)) {
            if (receipts) receiptStore.set(postId, receipts as ReceiptRecord[]);
        }
    }).catch(() => undefined);
}

function queueBatchFetch(postId: string): void {
    pendingPosts.add(postId);
    if (!batchTimer) batchTimer = window.setTimeout(flushBatch, 100);
}

function mountReceiptIndicator(mount: Element, postId: string): void {
    const render = () => {
        const vnode = React.createElement(ReceiptIndicator, {records: receiptStore.get(postId), isGroup: false});
        if ((ReactDOM as any).createRoot) {
            let root = roots.get(mount);
            if (!root) {
                root = (ReactDOM as any).createRoot(mount);
                roots.set(mount, root);
            }
            root.render(vnode);
        } else {
            ReactDOM.render(vnode, mount);
        }
    };
    receiptStore.subscribe(render);
    render();
}

function enhanceTimestamps(): void {
    document.querySelectorAll('.post__time:not([data-rr-mounted="true"])').forEach((timeEl) => {
        const postId = postIdFromTime(timeEl);
        if (!postId) return;
        (timeEl as HTMLElement).dataset.rrMounted = 'true';
        const mount = document.createElement('span');
        mount.className = 'rr-mount';
        timeEl.insertAdjacentElement('afterend', mount);
        mountReceiptIndicator(mount, postId);
        queueBatchFetch(postId);
    });
}

class ReadReceiptPlugin {
    public async initialize(registry: PluginRegistry): Promise<void> {
        injectCSS();
        registry.registerWebSocketEventHandler('custom_readreceipt_read_receipt_delivered', (msg: {data?: ReceiptRecord}) => { if (msg.data) receiptStore.upsert(msg.data); });
        registry.registerWebSocketEventHandler('custom_readreceipt_read_receipt_read', (msg: {data?: ReceiptRecord}) => { if (msg.data) receiptStore.upsert(msg.data); });
        stopObserver = startReadReceiptObserver();
        enhanceTimestamps();
        const mutation = new MutationObserver(enhanceTimestamps);
        mutation.observe(document.body, {childList: true, subtree: true});
    }

    public uninitialize(): void {
        stopObserver?.();
        if ((ReactDOM as any).createRoot) {
            roots.forEach((root) => root.unmount());
        } else {
            document.querySelectorAll('.rr-mount').forEach((el) => ReactDOM.unmountComponentAtNode(el));
        }
        roots.clear();
    }
}

declare global { interface Window { registerPlugin(id: string, plugin: ReadReceiptPlugin): void; } }

window.registerPlugin('readreceipt', new ReadReceiptPlugin());
