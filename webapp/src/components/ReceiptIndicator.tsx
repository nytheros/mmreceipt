import React from 'react';
import type {ReceiptRecord, ReceiptStatus} from '../store/receipt_store';

export function aggregateStatus(records: ReceiptRecord[]): ReceiptStatus {
    if (records.some((r) => r.status === 'read')) return 'read';
    if (records.some((r) => r.status === 'delivered')) return 'delivered';
    return 'sent';
}

export function tooltip(records: ReceiptRecord[], isGroup: boolean): string {
    const read = records.filter((r) => r.status === 'read');
    if (isGroup) return `Seen by ${read.length} users`;
    const first = read[0];
    if (!first?.updated_at) return 'Sent';
    return `Seen at ${new Date(first.updated_at).toLocaleTimeString([], {hour: '2-digit', minute: '2-digit'})}`;
}

export function ReceiptIndicator({records, isGroup = false}: {records: ReceiptRecord[]; isGroup?: boolean}) {
    const status = aggregateStatus(records);
    return <span className={`rr-indicator rr-${status}`} title={tooltip(records, isGroup)} style={{marginLeft: 4, fontWeight: 600}}>{status === 'sent' ? '✓' : '✓✓'}</span>;
}
