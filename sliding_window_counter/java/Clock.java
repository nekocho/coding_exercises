package java;

public interface Clock {
    long nowInSeconds();

    static final Real realInstance = new Real();

    public static Real real() {
        return realInstance;
    }

    public static Mock mock(final long atTime) {
        return new Mock(atTime);
    }

    public static Clock fixed() {
        return new Mock(System.currentTimeMillis() / 1000);
    }

    public static class Real implements Clock {
        @Override
        public long nowInSeconds() {
            return System.currentTimeMillis() / 1000;
        }
    }

    public static class Mock implements Clock {
        private long value;

        public Mock(long value) {
            this.value = value;
        }

        @Override
        public long nowInSeconds() {
            return value;
        }

        public void tick(final long seconds) {
            value += seconds;
        }
    }
}