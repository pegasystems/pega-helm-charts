#### JVM Arguments for better tunability
The following are the recommended arguments for better JVM tuning.

**— JVM CodeCache**<br>
CodeCache is the storage for JIT generated compiled code. If CodeCache is filled, compilation is done unless the code cache is flushed. <br>
**Flag**: `-XX:InitialCodeCacheSize=SIZE` <br>
**Purpose**: Sets the initial size of the JVM CodeCache. Set this to a reasonable size to avoid dynamic adjustment overhead.<br>
**JVM default**: 2MB <br>
**Pega Recommendation**:Set to 256M for large scale deployments and JVM default for others.

**Flag**: `-XX:ReservedCodeCacheSize=SIZE` <br>
**Purpose**: Sets the maximum size to which the JVM CodeCache can grow, after which the cache is attempted to be flushed. <br>
**JVM default**: 250MB<br>
**Pega Recommendation**:Set to 512M for large scale deployments and JVM default for others.<br> 

**— JVM Metaspace** <br>
Metaspace is a non-heap memory that stores the class metadata such as class structure, virtual method hierarchy, etc. 
This is the Java8+ equivalent of *PermGen* space. This space is unbounded by default, in the sense that it can endlessly expand to the system Native space.<br>
**Flag**: `-XX:MetaspaceSize=SIZE` <br>
**Purpose**: Sets the initial size of Metaspace at which point the JVM triggers Metaspace induced Garbage Collection (GC), in order to free Metaspace and prevent "Metaspace OutOfMemoryError".
 This is *not* the minimum Metaspace allocated. The setting allows you to set a higher Metaspace trigger size in order to avoid an initial GC that may not be warranted when the deployment starts up. <br>
 **JVM default**: 20MB <br>
**Recommendation**:Set to 512M for large scale deployments and JVM default for others.<br> 

**Flag**: `-XX:MaxMetaspaceSize=SIZE` <br>
**Purpose**: Sets the maximum value the JVM Metaspace can occupy. Beyond this "Metaspace OutOfMemoryError" is thrown.<br>
**JVM default**: 18446744073709486080 Bytes or Unbounded <br>
**Pega Recommendation**:  If left Unbounded (default) then there is a risk that the kernel will run out of physical memory, and when this occurs the system will hang without an error message, and therefore will be difficult to diagnose after recovery since there will be no logging detail to investigate.  Setting MaxMetaspaceSize to a minimum of 768m should be large enough to ensure that "Metaspace OutOfMemoryError" exceptions do not occur, and at the same time will prevent the kernel out of memory and resulting system hang without error messages.<br>

**Flag**: `-XX:+UseStringDeduplication` <br>
**Purpose**: Maintains a single copy of String literals *aka Deduplication*, in case multiple Strings have the same literal value<br>
**JVM default**: *UseStringDeduplication=false* equivalent to `-XX:-UseStringDeduplication`<br>
**Pega Recommendation**: Set to use the StringDeduplication.<br>

**Flag**: `-XX:+DisableExplicitGC` <br>
**Purpose**: Disables force GC, which can happen through external invocation of *System.gc()* call.<br>
**JVM default**: *DisableExplicitGC=false* equivalent to `-XX:-DisableExplicitGC`<br>
**Pega Recommendation**: Set to disable explicit GC<br>
**Note**: *This JVM argument is hardcoded in the Docker image with the recommended value. You cannot explicitly overwrite this argument.*

**Flag**: `-Djava.security.egd=file:///dev/urandom` <br>
**Purpose**: Allows generation of random numbers in non-blocking mode, in case the entropy is exhausted.<br>
**JVM default**: *dev/random*<br>
**Pega Recommendation**: Set to use *dev/urandom* <br>
**Note**: *This JVM argument is hardcoded in the Docker image with the recommended value. You cannot explicitly overwrite this argument.*

**Flag**: `-XX:+ExitOnOutOfMemoryError` <br>
**Purpose**: Exits the JVM when an OutOfMemoryError occurs. Exiting is required, as allowing the JVM to run after OutOfMemoryError leads to inconsistent system state.<br>
**JVM default**: *ExitOnOutOfMemoryError=false* equivalent to `-XX:-ExitOnOutOfMemoryError`<br>
**Pega Recommendation**: Set to Exit the JVM in case of OutOfMemoryError <br>
**Note**: *This JVM argument is hardcoded in the Docker image with the recommended value. You cannot explicitly overwrite this argument.*

**Flag**: `-Xlog:gc*,gc+heap=debug,gc+humongous=debug:file=/usr/local/tomcat/logs/gc.log:uptime,pid,level,time,tags:filecount=3,filesize=2M` <br>
**Purpose**: Dumps the GC logs with the provided tags into local log file. <br>
**JVM default**: No GC logs will be emitted by default.<br>
**Pega Recommendation**: Set to allow GC logs to be collected. Currently, the log rotation policy is set to 3 log files with maximum of 2MB each.<br>
<br>

**Flag**: `-Duser.timezone=TIMEZONE`<br>
**Purpose**: Sets the timezone the JVM will use.<br>
**JVM default**: By default, the JVM obtains time zone from the operating system.  The Pega docker image is configured to Etc/UTC.<br>
**Pega Recommendation**: Set the JVM time zone to match your database time zone. Set this parameter if your database time zone is not Etc/UTC.

**Flag**: `-XX:+HeapDumpOnOutOfMemoryError` <br>
**Purpose**: Generate a heap dump when an OutOfMemory Occurs.<br>
**JVM default**: JVM terminates the application when it encounters an OutOfMemoryError without generating heap dump<br>
**Pega Recommendation**: Set to generate heap dumps in case of OutOfMemoryError <br>