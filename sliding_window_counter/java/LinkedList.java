package java;

import java.util.LinkedList;

public class ListImpl implements Counter {
    private final Clock clock;
    private final long windowLengthInSeconds;
    private final LinkedList<Bucket> buckets = new LinkedList<Bucket>();
    private long count;

    public ListImpl(Clock clock, long windowLengthInSeconds) {
        this.clock = clock;
        this.windowLengthInSeconds = windowLengthInSeconds;
    }

    @Override
    public void increment() {
        final long t = clock.nowInSeconds();
        if (!buckets.isEmpty() && buckets.getLast().hasTs(t)) {
            buckets.getLast().increment();
        } else {
            final Bucket b = new Bucket(t, 1);
            buckets.addLast(b);
        }
        gc(t);
        count++;
    }

    @Override
    public long getValue() {
        gc(clock.nowInSeconds());
        return count;
    }

    private void gc(final long t) {
        final long cutOff = t-windowLengthInSeconds;
        while (!buckets.isEmpty() && buckets.getFirst().isEarlierThanOrEqualTo(cutOff)) {
            final Bucket discard = buckets.removeFirst();
            count -= discard.getValue();
        }
    }

    private static class Bucket implements Counter {
        private final long ts;
        private long value;

        public Bucket(long ts, long value) {
            this.ts = ts;
            this.value = value;
        }

        public boolean hasTs(final long value) {
            return ts == value;
        }

        public boolean isEarlierThanOrEqualTo(final long value) {
            return ts <= value;
        }

        @Override
        public void increment() {
            value++;
        }

        @Override
        public long getValue() {
            return value;
        }
    }
}
