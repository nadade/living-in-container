package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		host()
	case "container":
		container()
	default:
		panic("nope")
	}
}

func host() {
	fmt.Printf("[host] Running %v\n", os.Args[2:])

	cmd := exec.Command("/proc/self/exe", append([]string{"container"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:  syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
		Credential:  &syscall.Credential{Uid: 0, Gid: 0},
		UidMappings: []syscall.SysProcIDMap{{ContainerID: 0, HostID: 1000, Size: 1}},
		GidMappings: []syscall.SysProcIDMap{{ContainerID: 0, HostID: 1000, Size: 1}},
	}

	//create the new cgroup
	cg_path := "/sys/fs/cgroup/"
	my_cg_path := filepath.Join(cg_path, "container_gophercamp")
	must(os.MkdirAll(my_cg_path, 0755))

	//add this process to the new cgroup
	must(os.WriteFile(filepath.Join(my_cg_path, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))

	//limit max amount of processes to 20
	must(os.WriteFile(filepath.Join(my_cg_path, "pids.max"), []byte("20"), 0700))

	must(cmd.Run())
}

func container() {
	fmt.Printf("[container] Running %v\n", os.Args[2:])

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//set hostname within the container
	must(syscall.Sethostname([]byte("container")))

	//change root to the new file system (downloaded and untared in the vagrant provisioning shell)
	new_root := "/home/vagrant/ubuntufs"
	old_root := ".old_root"
	old_root_path := filepath.Join(new_root, old_root)

	//trickery - new root and old root must not be on the same file system
	must(syscall.Mount(new_root, new_root, "bind", syscall.MS_BIND|syscall.MS_REC, ""))
	must(os.MkdirAll(old_root_path, 0755))
	must(syscall.PivotRoot(new_root, old_root_path))

	//after changing root, we need to explicitly change directory to /
	must(os.Chdir("/"))

	//mount the procfs, otherwise /proc is empty
	must(syscall.Mount("proc", "/proc", "proc", 0, ""))

	//now we can get rid of the old root mount and empty dir
	old_root_path = filepath.Join("/", old_root)
	must(syscall.Unmount(old_root_path, syscall.MNT_DETACH))
	must(os.Remove(old_root_path))

	must(cmd.Run())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
