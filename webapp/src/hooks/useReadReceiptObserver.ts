const seen = new Set<string>();
const timers = new Map<string, number>();
const delivered = new Set<string>();

async function postReceipt(path: string, postId: string): Promise<void> {
    await fetch(`/plugins/readreceipt/api/v1/${path}`, {method: 'POST', credentials: 'same-origin', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({post_id: postId})});
}

function postIdFromElement(element: Element): string | undefined {
    const holder = element.closest('[id^="post_"]');
    return holder?.id.replace(/^post_/, '');
}

export function startReadReceiptObserver(): () => void {
    const observer = new IntersectionObserver((entries) => {
        for (const entry of entries) {
            const postId = postIdFromElement(entry.target);
            if (!postId || seen.has(postId)) continue;
            if (!delivered.has(postId)) { delivered.add(postId); void postReceipt('delivered', postId); }
            if (entry.intersectionRatio >= 0.7) {
                if (!timers.has(postId)) {
                    timers.set(postId, window.setTimeout(() => { seen.add(postId); timers.delete(postId); void postReceipt('read', postId); }, 2000));
                }
            } else {
                const timer = timers.get(postId);
                if (timer) window.clearTimeout(timer);
                timers.delete(postId);
            }
        }
    }, {threshold: [0, 0.7, 1]});

    const observe = () => document.querySelectorAll('.post:not([data-rr-observed="true"])').forEach((el) => { (el as HTMLElement).dataset.rrObserved = 'true'; observer.observe(el); });
    observe();
    const mutation = new MutationObserver(observe);
    mutation.observe(document.body, {childList: true, subtree: true});
    return () => { observer.disconnect(); mutation.disconnect(); timers.forEach((t) => window.clearTimeout(t)); };
}
