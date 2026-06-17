import React from 'react';
import {createRoot, Root} from 'react-dom/client';
type PluginRegistry = {registerWebSocketEventHandler(event: string, handler: (msg: {data?: ReceiptRecord}) => void): void};
import {ReceiptIndicator} from './components/ReceiptIndicator';
import {receiptStore, ReceiptRecord} from './store/receipt_store';
import {startReadReceiptObserver} from './hooks/useReadReceiptObserver';

const roots = new Map<Element, Root>();
let stopObserver: (() => void) | undefined;

function postIdFromTime(el: Element): string | undefined { return el.closest('[id^="post_"]')?.id.replace(/^post_/, ''); }

function enhanceTimestamps(): void {
    document.querySelectorAll('.post__time:not([data-rr-mounted="true"])').forEach((timeEl) => {
        const postId = postIdFromTime(timeEl);
        if (!postId) return;
        (timeEl as HTMLElement).dataset.rrMounted = 'true';
        const mount = document.createElement('span');
        mount.className = 'rr-mount';
        timeEl.insertAdjacentElement('afterend', mount);
        const root = createRoot(mount);
        roots.set(mount, root);
        const render = () => root.render(<ReceiptIndicator records={receiptStore.get(postId)} isGroup={false}/>);
        receiptStore.subscribe(render);
        render();
        void fetch(`/plugins/readreceipt/api/v1/status/${postId}`, {credentials: 'same-origin'}).then((r) => r.ok ? r.json() : undefined).then((data) => { if (data?.receipts) receiptStore.set(postId, data.receipts); }).catch(() => undefined);
    });
}

class ReadReceiptPlugin {
    public async initialize(registry: PluginRegistry): Promise<void> {
        registry.registerWebSocketEventHandler('custom_readreceipt_read_receipt_delivered', (msg: {data?: ReceiptRecord}) => { if (msg.data) receiptStore.upsert(msg.data); });
        registry.registerWebSocketEventHandler('custom_readreceipt_read_receipt_read', (msg: {data?: ReceiptRecord}) => { if (msg.data) receiptStore.upsert(msg.data); });
        stopObserver = startReadReceiptObserver();
        enhanceTimestamps();
        const mutation = new MutationObserver(enhanceTimestamps);
        mutation.observe(document.body, {childList: true, subtree: true});
    }

    public uninitialize(): void { stopObserver?.(); roots.forEach((root) => root.unmount()); roots.clear(); }
}

declare global { interface Window { registerPlugin(id: string, plugin: ReadReceiptPlugin): void; } }

window.registerPlugin('readreceipt', new ReadReceiptPlugin());
