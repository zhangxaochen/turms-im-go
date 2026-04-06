# Local Progress Tracker for Batch 2

[Context: ## resolveConflicts / closeSessionsWithConflictedDeviceTypes]
- [x] Java uses `SessionCloseStatus.DISCONNECTED_BY_CLIENT` to close conflicted sessions. Go uses `SessionCloseStatus_DISCONNECTED_BY_OTHER_DEVICE`, which is a different close status code.

---

[Context: ## resolveConflicts / closeSessionsWithConflictedDeviceTypes]
- [x] Java builds `nodeIdToDeviceTypes` by adding `deviceType` (the logging-in device type) to each node's set (line 817: `.add(deviceType)`). Go adds `conflictedDT` (the conflicted device type) to the node mapping. This means Java tells remote nodes to disconnect the logging-in device type, while Go tells them to disconnect the conflicted device types — fundamentally different conflict resolution behavior.

---

[Context: ## tryRegisterOnlineUser]
- [x] Java has a step (lines 656-676) that closes local sessions for device types that are already registered on Redis by other nodes. This handles edge cases like Redis crashes/restarts or lost connections to Redis. Go's `resolveConflicts` does not include this step — it only handles conflicted device types, not stale local sessions registered by other nodes.

---

[Context: ## tryRegisterOnlineUser]
- [x] Java has a session recovery path (lines 695-734) where if the session exists locally and its connection is closed, it returns the existing session for connection replacement (UDP-to-TCP/WebSocket recovery). This entire recovery branch is missing from Go.
- [x] Java calls `addOnlineDeviceIfAbsent` with `closeIdleSessionAfterSeconds`, `expectedNodeId`, and `expectedDeviceTimestamp` parameters for optimistic concurrency control. Go calls `AddOnlineDevice` without these parameters, losing the CAS-style protection against race conditions.

---

[Context: ## tryRegisterOnlineUser]
- [x] Java filters `deviceDetails` through `deviceDetailsItemPropertiesList` configuration before storing. Go passes `deviceDetails` directly without filtering.

---

[Context: ## tryRegisterOnlineUser]
- [x] Go calls `s.InvokeGoOnlineHandlers(ctx, nil, session)` passing `nil` as the `userSessionsManager`. Java calls `invokeGoOnlineHandlers(userSessionsManager, userSession)` with the actual non-null manager retrieved from `userIdToSessionsManager.computeIfAbsent`.

---

[Context: ## onSessionEstablished]
- [x] Java increments the `loggedInUsersCounter` metric. Go has an empty stub with a TODO comment.

---

[Context: ## onSessionEstablished]
- [x] Java conditionally pushes a session notification to the client via `userSessionsManager.pushSessionNotification(deviceType, serverId)` when `notifyClientsOfSessionInfoAfterConnected` is enabled. Go has an empty stub with no notification logic.

---

[Context: ## invokeGoOnlineHandlers]
- [x] Java invokes the `UserOnlineStatusChangeHandler.goOnline` plugin extension point via `pluginManager.invokeExtensionPointsSimultaneously`. Go has an empty stub with a TODO and does not invoke any plugin.

---

[Context: ## authAndUpdateHeartbeatTimestamp]
- [x] Java validates `deviceType` with `Validator.notNull` and `DeviceTypeUtil.validDeviceType`. Go does not validate the `deviceType` parameter at all.

---

[Context: ## notifySessionClosedListeners]
- [x] Java wraps each listener invocation in a try-catch block to prevent one failing listener from affecting others. Go calls listeners without any panic recovery, so a panic in one listener will propagate and abort the remaining listeners.
*Checked methods: getWsAddress(), getTcpAddress(), getUdpAddress()*
Now I have the full picture. Let me compare the three getter methods and the `updateCustomAddresses` logic.

---

[Context: ## getWsAddress()]
- [x] **Missing null/empty check when WS is disabled**: In Java, `wsAddress` is only updated inside `if (webSocketProperties.isEnabled())`, so it retains its previous value (initially `null`) when WS is disabled. In Go, `UpdateCustomAddresses` always unconditionally sets `m.wsAddress` regardless of whether WS is enabled or not. There is no `wsEnabled` parameter or check. The Go code always resolves and assigns the WS address even when it shouldn't.

---

[Context: ## getTcpAddress()]
- [x] **Missing enabled check when TCP is disabled**: Same issue as WS. Java only updates `tcpAddress` inside `if (tcpProperties.isEnabled())`. Go unconditionally sets `m.tcpAddress` every time `UpdateCustomAddresses` is called. There is no `tcpEnabled` parameter or conditional guard.

---

[Context: ## getUdpAddress()]
- [x] **Missing enabled check when UDP is disabled**: Same issue as WS/TCP. Java only updates `udpAddress` inside `if (udpProperties.isEnabled())`. Go unconditionally sets `m.udpAddress`. There is no `udpEnabled` parameter or conditional guard.
- [x] **Missing `BIND_ADDRESS` strategy handling**: Java's `queryHost` has a `BIND_ADDRESS` case that returns `bindHost` or errors if blank. Go falls through to a fallback `"127.0.0.1"` when `advertiseStrategy` is not `"ADVERTISE_ADDRESS"` and `host` is empty, which is incorrect behavior for `BIND_ADDRESS`.
- [x] **Missing `PRIVATE_ADDRESS` strategy handling**: Java queries the local private IP via `IpDetector.queryPrivateIp()`. Go has no equivalent — it falls through to the default `"127.0.0.1"` fallback.
- [x] **Missing `PUBLIC_ADDRESS` strategy handling**: Java queries the public IP via `IpDetector.queryPublicIp()`. Go has no equivalent — it falls through to the default `"127.0.0.1"` fallback.
- [x] **No error on blank advertiseHost for `ADVERTISE_ADDRESS`**: Java throws `IllegalArgumentException("The advertised host is not specified")` when `advertiseHost` is blank under `ADVERTISE_ADDRESS` strategy. Go silently falls through to check `host` and then `"127.0.0.1"`.
- [x] **Parameter name mismatch**: Java's second parameter is `bindHost` (the bind address), but Go names it `host`. This is cosmetic but contributes to confusion about what the parameter represents.
- [x] **Missing `areAddressPropertiesChange` optimization**: Java checks if discovery properties have actually changed before recomputing addresses, returning early with `FALSE` if they haven't. Go always recomputes unconditionally.
- [x] **Missing storage of `gatewayApiDiscoveryProperties`**: Java stores the current `DiscoveryProperties` in a field for later comparison in `areAddressPropertiesChange`. Go has no equivalent field or comparison logic.
- [x] **`queryHost` is synchronous in Go but asynchronous in Java**: Java's `queryHost` returns `Mono<String>` (reactive), and the address assignment happens in `.doOnNext()`. Go's `queryHost` is a synchronous function. This may cause issues if the Java version's async nature was relied upon (e.g., for `PUBLIC_ADDRESS` which does HTTP calls). The Go version cannot handle async IP detection.
*Checked methods: isConnected(), connect(), bind(boolean useFastBind, String dn, String password), search(String baseDn, Scope scope, DerefAliases derefAliases, int sizeLimit, int timeLimit, boolean typeOnly, List<String> attributes, String filter), modify(String dn, List<ModifyOperationChange> changes)*
Now I have all the files. Let me analyze each method carefully.
- **Java**: `return connection != null && !connection.isDisposed();`
- **Go**: `return c.Conn != nil && !c.Conn.IsClosing()`
This maps correctly. `IsClosing()` is the go-ldap equivalent of `isDisposed()`.
- **Java**: Uses an atomic CAS on `connectionMono` to ensure only one connection attempt at a time (singleton pattern). If `connectionMono` is already set, returns the existing Mono. The connection is stored in both `connectionMono` and `connection`.
- **Go**: Uses a mutex lock. Checks if connection is already established. Handles closing stale connections. Supports TLS.
The Go version changes semantics slightly: Java's `connect()` returns a `Mono` (lazy, reusable), while Go's returns `error` (eager). But functionally the connection-sharing logic is preserved. One notable difference: Java assigns to both `connectionMono` and `connection` fields, while Go only assigns to `c.Conn`. Since Go doesn't have the reactive stream separation, this is fine.
- **Java**: Sends a bind request with fast bind controls if `useFastBind` is true. Maps the response: success → true, `INVALID_CREDENTIALS` → false, other error codes → throws `LdapException`.
- **Go**: Ignores `useFastBind` entirely (comment says "We use Simple Bind regardless"). Maps `LDAPResultInvalidCredentials` → false, other errors → returns error.
**Bug**: The `useFastBind` parameter is accepted but ignored. In Java, when `useFastBind` is true, the request includes a `FAST_BIND` control (`ControlOidConst.FAST_BIND`). The Go version sends a plain bind without this control. This means the fast bind optimization (which tells the LDAP server to skip additional auth steps) is never applied.
- **Java**: Takes `baseDn, scope, derefAliases, sizeLimit, timeLimit, typeOnly, attributes, filter` and sends a `SearchRequest`, returning `Mono<SearchResult>`.
- **Go**: Takes the same parameters with `int` types for scope and derefAliases. Constructs a `ldap.NewSearchRequest` and returns the result.
This looks correct. The parameter mapping is faithful.
- **Java**: Returns `Mono.empty()` if changes is empty. Iterates changes: if type is `ADD` and attribute is empty, throws `LdapException`. Then sends the modify request.
- **Go**: Returns `nil` if changes is empty. Iterates changes: if type is `AddAttribute` and values is empty, returns an error. Then sends the modify request.
**Bug**: The error handling in Go differs from Java. In Java, the `LdapException` for empty ADD attributes is thrown immediately (synchronously) as an exception, which is the same as Go returning an error. However, in Java the loop checks **all** changes before sending the request — if a later change is invalid, the exception is still thrown before any request is sent. In Go, the validation and request building happen in the same loop, but the function also returns early on the first invalid change. This is actually functionally equivalent behavior (fail-fast on first invalid ADD).
Actually, let me re-read the Go code more carefully. The Go version checks `change.Type == ldap.AddAttribute && len(change.Values) == 0` and returns an error immediately (first invalid one). Java iterates all changes and throws on the first invalid one. The behavior is the same: fail on first invalid ADD attribute.
One more subtle difference: in Java, `modify` returns `Mono<ModifyResponse>` and the caller would need to check the response. In Go, `conn.Modify(req)` returns an error directly (which includes any LDAP error codes). This is a valid adaptation of the reactive pattern to the synchronous Go pattern.
---
Here are the bugs:
- [x] The `useFastBind` parameter is accepted but completely ignored. In Java, when `useFastBind` is true, the bind request includes a `FAST_BIND` control (`ControlOidConst.FAST_BIND` = a specific OID control with `criticality=false`). This control instructs the LDAP server to skip additional authentication steps (e.g., password policy evaluation). The Go version always performs a plain simple bind without this control, meaning the fast bind optimization is never applied regardless of the `useFastBind` argument value. [FIXED]
- [x] In Java, the validation for empty ADD attributes throws a specific `LdapException` with `ResultCode.INVALID_ATTRIBUTE_SYNTAX` and a descriptive message including the attribute type name. The Go version returns a generic `fmt.Errorf` with a string prefix "INVALID_ATTRIBUTE_SYNTAX:" but does not use a structured error type that callers can programmatically match against an LDAP result code. While this is arguably a stylistic difference for a Go codebase, it means callers cannot distinguish this error from other LDAP errors by result code as they could in Java. [FIXED]
*Checked methods: skipTag(), skipTagAndLength(), skipTagAndLengthAndValue(), readTag(), peekAndCheckTag(int tag), skipLength(), skipLengthAndValue(), writeLength(int length), readLength(), tryReadLengthIfReadable(), beginSequence(), beginSequence(int tag), endSequence(), writeBoolean(boolean value), writeBoolean(int tag, boolean value), readBoolean(), writeInteger(int value), writeInteger(int tag, int value), readInteger(), readIntWithTag(int tag), writeOctetString(String value), writeOctetString(byte[] value), writeOctetString(int tag, byte[] value), writeOctetString(byte[] value, int start, int length), writeOctetString(int tag, byte[] value, int start, int length), writeOctetString(int tag, String value), writeOctetStrings(List<String> values), readOctetString(), readOctetStringWithTag(int tag), readOctetStringWithLength(int length), writeEnumeration(int value), readEnumeration(), getBytes(), skipBytes(int length), close(), refCnt(), retain(), retain(int increment), touch(), touch(Object hint), release(), release(int decrement), isReadable(int length), isReadable(), isReadableWithEnd(int end), readerIndex()*
The `ttl_cache.go` file is not related to `BerBuffer` — it's a separate utility. I'll focus my review strictly on the `BerBuffer` port. Let me now do a detailed method-by-method comparison.

---

[Context: ## readLength()]
- [x] **Missing error handling for indefinite length (numBytes == 0)**: Java throws `DecodeException("Indefinite length is not supported")`, Go silently returns 0. This will cause the caller to misinterpret the structure.

---

[Context: ## readLength()]
- [x] **Missing validation for numBytes > 4**: Java throws `DecodeException("The length (...) is too long")`, Go silently proceeds with any numBytes value, potentially reading far beyond valid data.

---

[Context: ## readLength()]
- [x] **Missing insufficient data check**: Java checks `buffer.isReadable(lengthBytes)` before reading the length bytes and throws `DecodeException("Insufficient data")`. Go silently breaks out of the loop with a partial/truncated length value instead of throwing an error.

---

[Context: ## readLength()]
- [x] **Missing negative length check**: Java checks `if (length < 0)` after computation and throws `DecodeException("Invalid length bytes")`. Go returns the potentially corrupted length without validation.

---

[Context: ## tryReadLengthIfReadable()]
- [x] **Delegates to ReadLength() which has different error behavior**: Java's `tryReadLengthIfReadable` has its own full logic that returns -1 when there isn't enough data for the multi-byte length (line 155-157 in Java). The Go version simply delegates to `ReadLength()` which can silently return a truncated/incorrect value instead of -1 when multi-byte length data is partially available.

---

[Context: ## tryReadLengthIfReadable()]
- [x] **Missing indefinite length error**: Java throws `DecodeException("Indefinite length is not supported")` when numBytes == 0, but Go's `ReadLength()` returns 0 silently.

---

[Context: ## tryReadLengthIfReadable()]
- [x] **Missing numBytes > 4 error**: Java throws an exception, Go proceeds silently.

---

[Context: ## tryReadLengthIfReadable()]
- [x] **Missing negative length check**: Java validates `length < 0`, Go does not.

---

[Context: ## readBoolean()]
- [x] **Missing tag validation**: Java reads the tag and verifies it equals `TAG_BOOLEAN`, throwing `DecodeException` on mismatch. Go calls `SkipTagAndLength()` which blindly skips the tag without validating it.

---

[Context: ## readBoolean()]
- [x] **Missing length validation**: Java checks `length > 1` and throws `DecodeException("The boolean is too large")`, and checks `isReadable(length)`. Go skips the length without any validation.

---

[Context: ## readBoolean()]
- [x] **Returns false on insufficient data instead of throwing**: Java throws `DecodeException("Insufficient data")`. Go returns `false`, which is incorrect behavior — the caller cannot distinguish between "read a false boolean" and "failed to read".

---

[Context: ## readIntWithTag(int tag)]
- [x] **Returns 0 instead of throwing on tag mismatch**: Java throws `DecodeException` with tag mismatch details. Go silently returns 0, hiding protocol errors.

---

[Context: ## readIntWithTag(int tag)]
- [x] **Returns 0 instead of throwing when length > 4**: Java throws `DecodeException("The integer is too long")`. Go silently returns 0.

---

[Context: ## readIntWithTag(int tag)]
- [x] **Returns 0 instead of throwing on insufficient data**: Java throws `DecodeException("Insufficient data")`. Go returns 0.

---

[Context: ## readIntWithTag(int tag)]
- [x] **Incorrect negative value decoding logic**: Java extracts bits 0-6 of the first byte (`firstByte & 0x7F`) and then negates (`-value`). Go reads all bytes including the sign bit, then applies a sign-extension mask (`val |= ^((1 << (length * 8)) - 1)`). These produce different results for negative integers. For example, a 2-byte encoding of -129 (0xFF, 0x7F): Java computes `value = 0x7F = 127`, then `value = -127`. Go computes `val = 0xFF7F = 65407`, then `val |= ^0xFFFF = val | 0xFFFF0000 = 0xFFFFFF7F = -129`. The results differ (-127 vs -129).

---

[Context: ## readIntWithTag(int tag)]
- [x] **Returns 0 when length == 0**: Java would proceed to read a byte with `buffer.readByte()` and potentially throw. Go returns 0 silently.

---

[Context: ## writeInteger(int tag, int value)]
- [x] **Incorrect encoding for negative values**: Java's encoding uses bitmasks to determine the minimal encoding size for negative numbers (e.g., `(value & 0xFFFF_FF80) == 0xFFFF_FF80` checks if the value fits in 1 byte with sign extension). Go uses range checks (`value >= -128 && value <= 127`). While the range checks happen to produce the same byte count for negative values, the Go version writes different byte values. For example, for `value = -1`: Java writes `byte(value & 0xFF) = 0xFF`. Go writes `byte(-1) = 0xFF` — same in this case. But for `value = -129`: Java's 2-byte path writes `byte((value >> 8) & 0xFF) = 0xFF` and `byte(value & 0xFF) = 0x7F`. Go writes `byte(-129 >> 8) = 0xFF` and `byte(-129) = 0x7F`. These match. However, **the positive-value encoding differs**: Java masks the high byte with `& 0x7F` (e.g., `buffer.writeByte((value >> 8) & 0x7F)` for 2-byte positive values), while Go does not mask with `0x7F`. For positive values in the range 128-255, Java's 1-byte branch writes `value & 0x7F` (stripping bit 7), but Go's range check puts values > 127 into the 2-byte branch, so the behavior diverges. For value=200: Java's `(value & 0x0000_007F) == value` is `200 & 0x7F = 72 != 200`, so it falls to 2-byte: writes `byte((200 >> 8) & 0x7F) = 0` and `byte(200 & 0xFF) = 0xC8`. Go's `value >= -128 && value <= 127` is false for 200, so 2-byte: writes `byte(200 >> 8) = 0` and `byte(200) = 0xC8`. Same in this case. But for value=128: Java 2-byte path writes `byte((128>>8) & 0x7F) = 0` and `byte(128 & 0xFF) = 0x80`. Go writes same. So the positive paths are actually equivalent since Go's range boundaries align with Java's bitmask-based boundaries. This is correct.

---

[Context: ## writeOctetString(int tag, String value)]
- [x] **Missing length pre-reservation / backpatching logic**: Java reserves 3 bytes for the length (0x82 + 2 bytes), writes the UTF-8 string, then backpatches the actual length. Go writes the length immediately with `WriteLength(len(value))`, which uses `len(value)` (byte count of the Go string, which is already UTF-8). However, the critical difference is that Java always writes a 3-byte length (0x82 + 2 bytes) regardless of actual string length, while Go's `WriteLength` uses a variable-length encoding (1 byte for lengths <= 127, etc.). This means the wire format differs for short strings: Java produces 3 bytes of length overhead for all strings, Go produces 1 byte for strings <= 127 bytes. This is actually a more compact encoding, but it is not byte-compatible with the Java version. Whether this is a "bug" depends on whether wire-compatibility is required; for an LDAP client it may work fine since servers should handle both forms.

---

[Context: ## readOctetStringWithTag(int tag)]
- [x] **Returns empty string instead of throwing on tag mismatch**: Java throws `DecodeException` with tag details. Go silently returns `""`.

---

[Context: ## readOctetStringWithTag(int tag)]
- [x] **Returns empty string instead of throwing on insufficient data**: Java throws `DecodeException("Insufficient data")` when `!buffer.isReadable(length)`. Go's `ReadOctetStringWithLength` returns `""` silently.

---

[Context: ## readOctetStringWithLength(int length)]
- [x] **Returns empty string instead of throwing on insufficient data**: Java reads with `buffer.readCharSequence(length, UTF_8)` which would throw if insufficient data (Netty's ByteBuf does bounds checking). Go silently returns `""` when `b.readerIdx+length > len(b.buf)`.

---

[Context: ## beginSequence(int tag)]
- [x] **Different initialization of sequenceLengthWriterIndexes**: Java pre-allocates with `new int[8]` (fixed-size array with length 8) and tracks usage via `currentSequenceLengthIndex`. Go uses `append` to a dynamically grown slice starting from capacity 8 but length 0. This works but the slice index semantics differ — Go uses `append` which is fine.

---

[Context: ## beginSequence(int tag)]
- [x] **Missing overflow/growth check**: Java explicitly checks `currentSequenceLengthIndex >= writerIndexCount` and grows the array. Go relies on `append` to handle growth automatically. This is functionally equivalent.

---

[Context: ## endSequence()]
- [x] **Missing "Unbalanced sequences" validation**: Java throws `IllegalStateException("Unbalanced sequences")` when `--currentSequenceLengthIndex < 0`. Go silently does nothing (the `if b.currentSequenceLengthIndex > 0` check just skips).

---

[Context: ## endSequence()]
- [x] **Missing exception for length > 65535**: Java throws `DecodeException`. Go uses `panic`, which is a different error mechanism (not recoverable in normal Go flow without recover()).

---

[Context: ## peekAndCheckTag(int tag)]
- [x] **Minor difference — functionally equivalent but uses different comparison**: Java compares `buffer.getByte(buffer.readerIndex()) == tag` where `tag` is an `int` and `getByte` returns a `byte`. Go compares `int(b.buf[b.readerIdx]) == tag`. Both work correctly since Go's byte-to-int conversion matches Java's byte-to-int implicit widening (both produce unsigned 0-255 values). No bug here.

---

[Context: ## writeOctetString(String value)]
- [x] **Delegates to WriteOctetStringWithTag which uses different length encoding**: Same issue as `writeOctetString(int tag, String value)` — variable-length vs fixed 3-byte length encoding.

---

[Context: ## readTag()]
- [x] **Returns 0 on buffer underflow instead of throwing/panic**: Java's `buffer.readByte()` would throw `IndexOutOfBoundsException` if no bytes are readable. Go returns 0. This could be misinterpreted as a valid tag by callers.

---

[Context: ## skipTagAndLengthAndValue()]
- [x] **Functionally equivalent**: Correct.

---

[Context: ## skipLengthAndValue()]
- [x] **Functionally equivalent** (given ReadLength differences propagate).

---

[Context: ## isReadableWithEnd(int end)]
- [x] **Extra bounds check not present in Java**: Java's `isReadableWithEnd` only checks `buffer.readerIndex() < end`. Go adds `&& b.readerIdx < len(b.buf)`. The additional check means Go may return false when Java would return true (if `readerIndex < end` but `readerIdx >= len(buf)`). In practice this could cause early termination of parsing loops that rely on the Java semantics.

---

[Context: ## readBoolean() (second look)]
- [x] **Reads only 1 byte of value, ignoring length field**: Java reads the length via `readLength()` and validates it, but then reads exactly 1 byte. Go calls `SkipTagAndLength()` which consumes the tag byte and length bytes, then reads 1 byte. If the length indicated more than 1 byte of content, Go would leave the reader index misaligned compared to Java which also only reads 1 byte but validates the length first. The Java code validates `length > 1` and throws an error; Go silently ignores extra bytes.
---
Here is the consolidated bug list:

---

[Context: ## readLength]
- [x] Silently returns 0 for indefinite length (numBytes == 0) instead of throwing an error like Java's `DecodeException("Indefinite length is not supported")`

---

[Context: ## readLength]
- [x] Missing validation that numBytes <= 4; Java throws `DecodeException` for lengths > 4, Go silently reads up to any number of bytes

---

[Context: ## readLength]
- [x] Silently breaks with a truncated/partial length when insufficient bytes are available instead of throwing `DecodeException("Insufficient data")`

---

[Context: ## readLength]
- [x] Missing negative length validation; Java checks `length < 0` and throws `DecodeException("Invalid length bytes")`, Go does not