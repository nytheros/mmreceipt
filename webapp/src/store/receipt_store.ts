export type ReceiptStatus = 'sent' | 'delivered' | 'read';
export type ReceiptRecord = {post_id: string; user_id: string; status: ReceiptStatus; updated_at?: number};

type Listener = () => void;

class ReceiptStore {
    private receipts = new Map<string, ReceiptRecord[]>();
    private listeners = new Set<Listener>();

    subscribe(listener: Listener): () => void { this.listeners.add(listener); return () => this.listeners.delete(listener); }
    emit(): void { for (const listener of this.listeners) listener(); }
    get(postId: string): ReceiptRecord[] { return this.receipts.get(postId) || []; }
    set(postId: string, records: ReceiptRecord[]): void { this.receipts.set(postId, records); this.emit(); }
    upsert(record: ReceiptRecord): void {
        const current = this.get(record.post_id).filter((r) => r.user_id !== record.user_id);
        current.push(record);
        this.receipts.set(record.post_id, current);
        this.emit();
    }
}

export const receiptStore = new ReceiptStore();
