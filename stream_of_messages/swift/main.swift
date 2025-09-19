class MessageStream {
    enum Endianness {
        case little
        case big
    }

    enum State {
        case waitingForLength
        case waitingForMessage(UInt32)
    }

    let endianness: Endianness = .big
    var state: State = .waitingForLength
    var buffer: [UInt8] = []

    private func decodeLengthByteShifts() -> UInt32 {
        switch self.endianness {
        case .little:
            return UInt32(self.buffer[0]) + (UInt32(self.buffer[1]) << 8) + (UInt32(self.buffer[2]) << 16) + (UInt32(self.buffer[3]) << 24)
        case .big:
            return UInt32(self.buffer[3]) + (UInt32(self.buffer[2]) << 8) + (UInt32(self.buffer[1]) << 16) + (UInt32(self.buffer[0]) << 24)
        }
    }

    private func decodeLength() -> UInt32 {
        var value: UInt32 = 0
        self.buffer.withUnsafeBytes { bufferPtr in
            withUnsafeMutableBytes(of: &value) { valuePtr in
                valuePtr.copyMemory(from: bufferPtr)
            }
        }
        switch self.endianness {
        case .little:
            return value.littleEndian
        case .big:
            return value.bigEndian
        }
    }

    func accept(chunk: [UInt8]) {
        self.buffer.append(contentsOf: chunk)
        switch self.state {
        case .waitingForLength:
            if self.buffer.count >= 4 {
                let length = self.decodeLength()
                self.state = .waitingForMessage(length)
                self.buffer = Array(self.buffer[4...])
            }
        case .waitingForMessage(let length):
            if self.buffer.count >= length {
                self.processMessage(message: Array(self.buffer.prefix(Int(length))))
                self.state = .waitingForLength
                self.buffer = Array(self.buffer[Int(length)...])
            }
        }
    }

    func process(message: [UInt8]) {
        print(message.map { $0 })
    }
}
