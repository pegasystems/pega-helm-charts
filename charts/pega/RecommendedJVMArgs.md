#### JVM Arguments for better tunability
The following are the recommended arguments for better JVM tuning.

**— JVM CodeCache**<br>
CodeCache is the storage for JIT generated compiled code. If CodeCache is filled, compilation stops unless the code cache is flushed. <br>
**Flag**: `-XX:+UseCodeCacheFlushing` <br>
**Purpose**: Attempts to clean the CodeCache when filled, before stopping the compilation<br>
**Defaults and Recommendation**: JVM default<br> 

**Flag**: `-XX:InitialCodeCacheSize=SIZE` <br>
**Purpose**: The initial value of the JVM CodeCache. Set this value to reasonable size to avoid overhead in dynamically increasing the cache.<br>
**Defaults and Recommendation**:Set to `256M` for large scale deployments and JVM defaults for others.<br> 

**Flag**: `-XX:+ReservedCodeCacheSize=SIZE` <br>
**Purpose**: The maximum value the JVM CodeCache can grow, after which the cache is attempted to be flushed. <br>
**Defaults and Recommendation**:Set to `512M` for large scale deployments and JVM defaults for others.<br> 

**— JVM Metaspace** <br>
Metaspace is a non-heap memory that stores the class metadata such as class structure, virtual method hierarchy, etc. 
This is the Java8+ equivalent of *PermGen* space. This space is unbounded by default, in the sense that it can endlessly expand to the system Native space.<br>
**Flag**: `-XX:MetaspaceSize=SIZE` <br>
**Purpose**: This defines the initial size of Metaspace at which GC should be induced, expecting to clear up some Metaspace.
 Note that this is *not* the initial size the metaspace is allocated. It is set just to avoid the early 'Metaspace induced GC'.<br>
**Recommendation**: `512M` for large scale deployments and defaults for others.<br> 

**Flag**: `-XX:MaxMetaspaceSize=SIZE` <br>
**Purpose**: This is the maximum value the JVM Metaspace occupies. Beyond this "Metaspace OutOfMemoryError" is thrown.<br>
**Defaults and Recommendation**: Unbounded, which is the JVM default value.<br>

**Flag**: `-XX:+UseStringDeduplication` <br>
**Purpose**: When multiple Java Strings have same literal value, these strings are duplicated. This flag deals with optimizing these duplications, maintaining a single String.<br>
**Defaults and Recommendation**: Set to use the string deduplication. Tune this according to the needs.<br>

**Flag**: `-XX:+DisableExplicitGC` <br>
**Purpose**: Disables force GC, which can happen through external invocation of *System.gc()* calls.<br>
**Defaults and Recommendation**: Set to disable explicit GC<br>
**Note**: *This flag cannot be overridden through explicitly passing the flag*.

**Flag**: `-Djava.security.egd=file:///dev/urandom` <br>
**Purpose**: For non-bocking random number generation in case the entropy is exhausted.<br>
**Defaults and Recommendation**: Set to use *urandom* <br>
**Note**: *This flag cannot be overridden through explicitly passing the flag*.

**Flag**: `-XX:+CrashOnOutOfMemoryError` <br>
**Purpose**: When an OOME is occured, the thread in which the error occurs is terminated(if the error is not caught).
  However, this leaves the JVM still running and in an inconsistent state. In this case, it's better to crash the JVM.<br>
**Defaults and Recommendation**: Set to Crash the JVM in case of OOME<br>
**Note**: *This flag cannot be overridden through explicitly passing the flag*.

**Flag**: `-Xlog:gc*,gc+heap=debug,gc+humongous=debug:file=/usr/local/tomcat/logs/gc.log:uptime,pid,level,time,tags:filecount=3,filesize=2M` <br>
**Purpose**: Dumps the GC logs with the provided tags into local log file. <br>
**Recommendation**: Set to allow GC logs to be collected. Currently, the log rotation policy is set to 3 log files with maximum of 2MB each.<br>
<br>
##### Few more arguments:
**Flag**: `-XX:MaxGCPauseMillis=n`<br>
**Purpose**: Sets the peak pause time for GC to happen.<br>
**Recommendation**: Set to higher value for BATCH pods for better throughput.
<br>

**Flag**: `-XX:StringTableSize=n`<br>
**Purpose**: Sets the number of hash-buckets to be used in string pool for accommodating Strings<br>
**Recommendation**: JVM default is 60013. Tune this based on requirements.<br>







