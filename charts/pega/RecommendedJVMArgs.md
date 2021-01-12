#### JVM Arguments for better tunability
The following are the recommended arguments for better JVM tuning.

**— JVM CodeCache**<br>
CodeCache is the storage for JIT generated compiled code. If CodeCache is filled, compilation is done unless the code cache is flushed. <br>
**Flag**: `-XX:InitialCodeCacheSize=SIZE` <br>
**Purpose**: Sets the initial size of the JVM CodeCache. Set this to a reasonable size to avoid dynamic adjustment overhead.<br>
**JVM default**: 2MB <br>
**Pega Recommendation**:Set to 256M for large scale deployments and JVM default for others.

**Flag**: `-XX:ReservedCodeCacheSize=SIZE` <br>
**Purpose**: Sets the maximum size upto which the JVM CodeCache can grow, after which the cache is attempted to be flushed. <br>
**JVM default**: 250MB<br>
**Pega Recommendation**:Set to 512M for large scale deployments and JVM default for others.<br> 

**— JVM Metaspace** <br>
Metaspace is a non-heap memory that stores the class metadata such as class structure, virtual method hierarchy, etc. 
This is the Java8+ equivalent of *PermGen* space. This space is unbounded by default, in the sense that it can endlessly expand to the system Native space.<br>
**Flag**: `-XX:MetaspaceSize=SIZE` <br>
**Purpose**: Sets the initial size of Metaspace at which GC should be induced, expecting to clear up some Metaspace.
 Note that this is *not* the initial size the metaspace is allocated. It is set just to avoid the early Metaspace induced GC.<br>
 **JVM default**: 20MB <br>
**Recommendation**:Set to 512M for large scale deployments and JVM default for others.<br> 

**Flag**: `-XX:MaxMetaspaceSize=SIZE` <br>
**Purpose**: Sets the maximum value the JVM Metaspace can occupy. Beyond this "Metaspace OutOfMemoryError" is thrown.<br>
**JVM default**: 18446744073709486080 Bytes or Unbounded <br>
**Pega Recommendation**: Unbounded, which is the JVM default value.<br>

**Flag**: `-XX:+UseStringDeduplication` <br>
**Purpose**: Maintains a single copy of String literals *aka Deduplication*, in case multiple Strings having same literal value<br>
**JVM default**: *UseStringDeduplication=false* equivalent to `-XX:-UseStringDeduplication`<br>
**Pega Recommendation**: Set to use the StringDeduplication.<br>

**Flag**: `-XX:+DisableExplicitGC` <br>
**Purpose**: Disables force GC, which can happen through external invocation of *System.gc()* call.<br>
**JVM default**: *DisableExplicitGC=false* equivalent to `-XX:-DisableExplicitGC`<br>
**Pega Recommendation**: Set to disable explicit GC<br>
**Note**: *This flag cannot be overridden through explicitly passing the flag*.

**Flag**: `-Djava.security.egd=file:///dev/urandom` <br>
**Purpose**: Allows generation of random numbers in non-bocking mode, in case the entropy is exhausted.<br>
**JVM default**: *dev/random*<br>
**Pega Recommendation**: Set to use *dev/urandom* <br>
**Note**: *This flag cannot be overridden through explicitly passing the flag*.

**Flag**: `-XX:+ExitOnOutOfMemoryError` <br>
**Purpose**: Exits the JVM when an OutOfMemoryError occurs. Exiting is required, as allowing the JVM to run after OOME leads to inconsistent system state.<br>
**JVM default**: *ExitOnOutOfMemoryError=false* equivalent to `-XX:-ExitOnOutOfMemoryError`<br>
**Pega Recommendation**: Set to Exit the JVM in case of OOME <br>
**Note**: *This flag cannot be overridden through explicitly passing the flag*.

**Flag**: `-Xlog:gc*,gc+heap=debug,gc+humongous=debug:file=/usr/local/tomcat/logs/gc.log:uptime,pid,level,time,tags:filecount=3,filesize=2M` <br>
**Purpose**: Dumps the GC logs with the provided tags into local log file. <br>
**JVM default**: No GC logs will be emitted by default.<br>
**Pega Recommendation**: Set to allow GC logs to be collected. Currently, the log rotation policy is set to 3 log files with maximum of 2MB each.<br>
<br>
##### Few more arguments:
**Flag**: `-XX:MaxGCPauseMillis=n`<br>
**Purpose**: Sets the peak pause time for GC to happen.<br>
**JVM default**: 200 milliseconds<br>
**Pega Recommendation**: Set to higher value for BATCH pods for better throughput.
<br>

**Flag**: `-XX:StringTableSize=n`<br>
**Purpose**: Sets the number of hash-buckets to be used in string pool for accommodating Strings<br>
**JVM default**: 65536, the default Hash bucket size of JVM String pool <br>
**Pega Recommendation**: Tune this based on requirements.<br>







