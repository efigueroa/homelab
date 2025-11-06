# NVIDIA GPU Acceleration Setup (GTX 1070)

This guide covers setting up NVIDIA GPU acceleration for your homelab running on **Proxmox 9 (Debian 13)** with an **NVIDIA GTX 1070**.

## Overview

GPU acceleration provides significant benefits:
- **Jellyfin**: Hardware video transcoding (H.264, HEVC)
- **Immich**: Faster ML inference (face recognition, object detection)
- **Performance**: 10-20x faster transcoding vs CPU
- **Efficiency**: Lower power consumption, CPU freed for other tasks

**Your Hardware:**
- **GPU**: NVIDIA GTX 1070 (Pascal architecture)
- **Capabilities**: NVENC (encoding), NVDEC (decoding), CUDA
- **Max Concurrent Streams**: 2 (can be unlocked)
- **Supported Codecs**: H.264, HEVC (H.265)

## Architecture Overview

```
Proxmox Host (Debian 13)
  │
  ├─ NVIDIA Drivers (host)
  ├─ NVIDIA Container Toolkit
  │
  └─ Docker VM/LXC
       │
       ├─ GPU passthrough
       │
       └─ Jellyfin/Immich containers
            └─ Hardware transcoding
```

## Part 1: Proxmox Host Setup

### Step 1.1: Enable IOMMU (for GPU Passthrough)

**Edit GRUB configuration:**

```bash
# SSH into Proxmox host
ssh root@proxmox-host

# Edit GRUB config
nano /etc/default/grub
```

**Find this line:**
```
GRUB_CMDLINE_LINUX_DEFAULT="quiet"
```

**Replace with (Intel CPU):**
```
GRUB_CMDLINE_LINUX_DEFAULT="quiet intel_iommu=on iommu=pt"
```

**Or (AMD CPU):**
```
GRUB_CMDLINE_LINUX_DEFAULT="quiet amd_iommu=on iommu=pt"
```

**Update GRUB and reboot:**
```bash
update-grub
reboot
```

**Verify IOMMU is enabled:**
```bash
dmesg | grep -e DMAR -e IOMMU

# Should see: "IOMMU enabled"
```

### Step 1.2: Load VFIO Modules

**Edit modules:**
```bash
nano /etc/modules
```

**Add these lines:**
```
vfio
vfio_iommu_type1
vfio_pci
vfio_virqfd
```

**Update initramfs:**
```bash
update-initramfs -u -k all
reboot
```

### Step 1.3: Find GPU PCI ID

```bash
lspci -nn | grep -i nvidia

# Example output:
# 01:00.0 VGA compatible controller [0300]: NVIDIA Corporation GP104 [GeForce GTX 1070] [10de:1b81] (rev a1)
# 01:00.1 Audio device [0403]: NVIDIA Corporation GP104 High Definition Audio Controller [10de:10f0] (rev a1)
```

**Note the IDs**: `10de:1b81` and `10de:10f0` (your values may differ)

### Step 1.4: Configure VFIO

**Create VFIO config:**
```bash
nano /etc/modprobe.d/vfio.conf
```

**Add (replace with your IDs from above):**
```
options vfio-pci ids=10de:1b81,10de:10f0
softdep nvidia pre: vfio-pci
```

**Blacklist nouveau (open-source NVIDIA driver):**
```bash
echo "blacklist nouveau" >> /etc/modprobe.d/blacklist.conf
```

**Update and reboot:**
```bash
update-initramfs -u -k all
reboot
```

**Verify GPU is bound to VFIO:**
```bash
lspci -nnk -d 10de:1b81

# Should show:
# Kernel driver in use: vfio-pci
```

## Part 2: VM/LXC Setup

### Option A: Using VM (Recommended for Docker)

**Create Ubuntu 24.04 VM with GPU passthrough:**

1. **Create VM in Proxmox UI**:
   - OS: Ubuntu 24.04 Server
   - CPU: 4+ cores
   - RAM: 16GB+
   - Disk: 100GB+

2. **Add PCI Device** (GPU):
   - Hardware → Add → PCI Device
   - Device: Select your GTX 1070 (01:00.0)
   - ✅ All Functions
   - ✅ Primary GPU (if no other GPU)
   - ✅ PCI-Express

3. **Add PCI Device** (GPU Audio):
   - Hardware → Add → PCI Device
   - Device: NVIDIA Audio (01:00.1)
   - ✅ All Functions

4. **Machine Settings**:
   - Machine: q35
   - BIOS: OVMF (UEFI)
   - Add EFI Disk

5. **Start VM** and install Ubuntu

### Option B: Using LXC (Advanced, Less Stable)

**Note**: LXC with GPU is less reliable. VM recommended.

If you insist on LXC:
```bash
# Edit LXC config
nano /etc/pve/lxc/VMID.conf

# Add:
lxc.cgroup2.devices.allow: c 195:* rwm
lxc.cgroup2.devices.allow: c 509:* rwm
lxc.mount.entry: /dev/nvidia0 dev/nvidia0 none bind,optional,create=file
lxc.mount.entry: /dev/nvidiactl dev/nvidiactl none bind,optional,create=file
lxc.mount.entry: /dev/nvidia-uvm dev/nvidia-uvm none bind,optional,create=file
```

**For this guide, we'll use VM (Option A)**.

## Part 3: VM Guest Setup (Debian 13)

Now we're inside the Ubuntu/Debian VM where Docker runs.

### Step 3.1: Install NVIDIA Drivers

**SSH into your Docker VM:**
```bash
ssh user@docker-vm
```

**Update system:**
```bash
sudo apt update
sudo apt upgrade -y
```

**Debian 13 - Install NVIDIA drivers:**
```bash
# Add non-free repositories
sudo nano /etc/apt/sources.list

# Add 'non-free non-free-firmware' to each line, example:
deb http://deb.debian.org/debian bookworm main non-free non-free-firmware
deb http://deb.debian.org/debian bookworm-updates main non-free non-free-firmware

# Update and install
sudo apt update
sudo apt install -y linux-headers-$(uname -r)
sudo apt install -y nvidia-driver nvidia-smi

# Reboot
sudo reboot
```

**Verify driver installation:**
```bash
nvidia-smi

# Should show:
# +-----------------------------------------------------------------------------+
# | NVIDIA-SMI 535.xx.xx    Driver Version: 535.xx.xx    CUDA Version: 12.2     |
# |-------------------------------+----------------------+----------------------+
# | GPU  Name        Persistence-M| Bus-Id        Disp.A | Volatile Uncorr. ECC |
# | Fan  Temp  Perf  Pwr:Usage/Cap|         Memory-Usage | GPU-Util  Compute M. |
# |===============================+======================+======================|
# |   0  NVIDIA GeForce ...  Off  | 00000000:01:00.0 Off |                  N/A |
# | 30%   35C    P8    10W / 150W |      0MiB /  8192MiB |      0%      Default |
# +-------------------------------+----------------------+----------------------+
```

✅ **Success!** Your GTX 1070 is now accessible in the VM.

### Step 3.2: Install NVIDIA Container Toolkit

**Add NVIDIA Container Toolkit repository:**
```bash
curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey | sudo gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg

curl -s -L https://nvidia.github.io/libnvidia-container/stable/deb/nvidia-container-toolkit.list | \
  sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#g' | \
  sudo tee /etc/apt/sources.list.d/nvidia-container-toolkit.list
```

**Install toolkit:**
```bash
sudo apt update
sudo apt install -y nvidia-container-toolkit
```

**Configure Docker to use NVIDIA runtime:**
```bash
sudo nvidia-ctk runtime configure --runtime=docker
```

**Restart Docker:**
```bash
sudo systemctl restart docker
```

**Verify Docker can access GPU:**
```bash
docker run --rm --gpus all nvidia/cuda:12.2.0-base-ubuntu22.04 nvidia-smi

# Should show nvidia-smi output from inside container
```

✅ **Success!** Docker can now use your GPU.

## Part 4: Configure Jellyfin for GPU Transcoding

### Step 4.1: Update Jellyfin Compose File

**Edit compose file:**
```bash
cd ~/homelab/compose/media/frontend/jellyfin
nano compose.yaml
```

**Uncomment the GPU sections:**

```yaml
services:
  jellyfin:
    container_name: jellyfin
    image: lscr.io/linuxserver/jellyfin:latest
    env_file:
      - .env
    volumes:
      - ./config:/config
      - ./cache:/cache
      - /media/movies:/media/movies:ro
      - /media/tv:/media/tv:ro
      - /media/music:/media/music:ro
      - /media/photos:/media/photos:ro
      - /media/homemovies:/media/homemovies:ro
    ports:
      - "8096:8096"
      - "7359:7359/udp"
    restart: unless-stopped
    networks:
      - homelab
    labels:
      traefik.enable: true
      traefik.http.routers.jellyfin.rule: Host(`flix.fig.systems`) || Host(`flix.edfig.dev`)
      traefik.http.routers.jellyfin.entrypoints: websecure
      traefik.http.routers.jellyfin.tls.certresolver: letsencrypt
      traefik.http.services.jellyfin.loadbalancer.server.port: 8096

    # UNCOMMENT THESE LINES FOR GTX 1070:
    runtime: nvidia
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: all
              capabilities: [gpu]

networks:
  homelab:
    external: true
```

**Restart Jellyfin:**
```bash
docker compose down
docker compose up -d
```

**Check logs:**
```bash
docker compose logs -f

# Should see lines about NVENC/CUDA being detected
```

### Step 4.2: Enable in Jellyfin UI

1. Go to https://flix.fig.systems
2. Dashboard → Playback → Transcoding
3. **Hardware acceleration**: NVIDIA NVENC
4. **Enable hardware decoding for**:
   - ✅ H264
   - ✅ HEVC
   - ✅ VC1
   - ✅ VP8
   - ✅ MPEG2
5. **Enable hardware encoding**
6. **Enable encoding in HEVC format**
7. Save

### Step 4.3: Test Transcoding

1. Play a video in Jellyfin web UI
2. Click Settings (gear icon) → Quality
3. Select a lower bitrate to force transcoding
4. In another terminal:
   ```bash
   nvidia-smi

   # While video is transcoding, should see:
   # GPU utilization: 20-40%
   # Memory usage: 500-1000MB
   ```

✅ **Success!** Jellyfin is using your GTX 1070!

## Part 5: Configure Immich for GPU Acceleration

Immich can use GPU for two purposes:
1. **ML Inference** (face recognition, object detection)
2. **Video Transcoding**

### Step 5.1: ML Inference (CUDA)

**Edit Immich compose file:**
```bash
cd ~/homelab/compose/media/frontend/immich
nano compose.yaml
```

**Change ML image to CUDA version:**

Find this line:
```yaml
image: ghcr.io/immich-app/immich-machine-learning:${IMMICH_VERSION:-release}
```

Change to:
```yaml
image: ghcr.io/immich-app/immich-machine-learning:${IMMICH_VERSION:-release}-cuda
```

**Add GPU support:**

```yaml
  immich-machine-learning:
    container_name: immich_machine_learning
    image: ghcr.io/immich-app/immich-machine-learning:${IMMICH_VERSION:-release}-cuda
    volumes:
      - model-cache:/cache
    env_file:
      - .env
    restart: always
    networks:
      - immich_internal

    # ADD THESE LINES:
    runtime: nvidia
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: all
              capabilities: [gpu]
```

### Step 5.2: Video Transcoding (NVENC)

**For video transcoding, add to immich-server:**

```yaml
  immich-server:
    container_name: immich_server
    image: ghcr.io/immich-app/immich-server:${IMMICH_VERSION:-release}
    # ... existing config ...

    # ADD THESE LINES:
    runtime: nvidia
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: all
              capabilities: [gpu]
```

**Restart Immich:**
```bash
docker compose down
docker compose up -d
```

### Step 5.3: Enable in Immich UI

1. Go to https://photos.fig.systems
2. Administration → Settings → Video Transcoding
3. **Transcoding**: h264 (NVENC)
4. **Hardware Acceleration**: NVIDIA
5. Save

6. Administration → Settings → Machine Learning
7. **Facial Recognition**: Enabled
8. **Object Detection**: Enabled
9. Should automatically use CUDA

### Step 5.4: Test ML Inference

1. Upload photos with faces
2. In terminal:
   ```bash
   nvidia-smi

   # While processing, should see:
   # GPU utilization: 50-80%
   # Memory usage: 2-4GB
   ```

✅ **Success!** Immich is using GPU for ML inference!

## Part 6: Performance Tuning

### GTX 1070 Specific Settings

**Jellyfin optimal settings:**
- Hardware acceleration: NVIDIA NVENC
- Target transcode bandwidth: Let clients decide
- Enable hardware encoding: Yes
- Prefer OS native DXVA or VA-API hardware decoders: No
- Allow encoding in HEVC format: Yes (GTX 1070 supports HEVC)

**Immich optimal settings:**
- Transcoding: h264 or hevc
- Target resolution: 1080p (for GTX 1070)
- CRF: 23 (good balance)
- Preset: fast

### Unlock NVENC Stream Limit

GTX 1070 is limited to 2 concurrent transcoding streams. You can unlock unlimited streams:

**Install patch:**
```bash
# Inside Docker VM
git clone https://github.com/keylase/nvidia-patch.git
cd nvidia-patch
sudo bash ./patch.sh

# Reboot
sudo reboot
```

**Verify:**
```bash
nvidia-smi

# Now supports unlimited concurrent streams
```

⚠️ **Note**: This is a hack that modifies NVIDIA driver. Use at your own risk.

### Monitor GPU Usage

**Real-time monitoring:**
```bash
watch -n 1 nvidia-smi
```

**Check GPU usage from Docker:**
```bash
docker stats $(docker ps --format '{{.Names}}' | grep -E 'jellyfin|immich')
```

## Troubleshooting

### GPU Not Detected in VM

**Check from Proxmox host:**
```bash
lspci | grep -i nvidia
```

**Check from VM:**
```bash
lspci | grep -i nvidia
nvidia-smi
```

**If not visible in VM:**
1. Verify IOMMU is enabled (`dmesg | grep IOMMU`)
2. Check PCI passthrough is configured correctly
3. Ensure VM is using q35 machine type
4. Verify BIOS is OVMF (UEFI)

### Docker Can't Access GPU

**Error**: `could not select device driver "" with capabilities: [[gpu]]`

**Fix:**
```bash
# Reconfigure NVIDIA runtime
sudo nvidia-ctk runtime configure --runtime=docker
sudo systemctl restart docker

# Test again
docker run --rm --gpus all nvidia/cuda:12.2.0-base-ubuntu22.04 nvidia-smi
```

### Jellyfin Shows "No Hardware Acceleration Available"

**Check:**
```bash
# Verify container has GPU access
docker exec jellyfin nvidia-smi

# Check Jellyfin logs
docker logs jellyfin | grep -i nvenc
```

**Fix:**
1. Ensure `runtime: nvidia` is uncommented
2. Verify `deploy.resources.reservations.devices` is configured
3. Restart container: `docker compose up -d`

### Transcoding Fails with "Failed to Open GPU"

**Check:**
```bash
# GPU might be busy
nvidia-smi

# Kill processes using GPU
sudo fuser -v /dev/nvidia*
```

### Low GPU Utilization During Transcoding

**Normal**: GTX 1070 is powerful. 20-40% utilization is expected for single stream.

**To max out GPU:**
- Transcode multiple streams simultaneously
- Use higher resolution source (4K)
- Enable HEVC encoding

## Performance Benchmarks (GTX 1070)

**Typical Performance:**
- **4K HEVC → 1080p H.264**: ~120-150 FPS (real-time)
- **1080p H.264 → 720p H.264**: ~300-400 FPS
- **Concurrent streams**: 4-6 (after unlocking limit)
- **Power draw**: 80-120W during transcoding
- **Temperature**: 55-65°C

**Compare to CPU (typical 4-core):**
- **4K HEVC → 1080p H.264**: ~10-15 FPS
- CPU would be at 100% utilization
- GPU: 10-15x faster!

## Monitoring and Maintenance

### Create GPU Monitoring Dashboard

**Install nvtop (nvidia-top):**
```bash
sudo apt install nvtop
```

**Run:**
```bash
nvtop
```

Shows real-time GPU usage, memory, temperature, processes.

### Check GPU Health

```bash
# Temperature
nvidia-smi --query-gpu=temperature.gpu --format=csv

# Memory usage
nvidia-smi --query-gpu=memory.used,memory.total --format=csv

# Fan speed
nvidia-smi --query-gpu=fan.speed --format=csv

# Power draw
nvidia-smi --query-gpu=power.draw,power.limit --format=csv
```

### Automated Monitoring

Add to cron:
```bash
crontab -e

# Add:
*/5 * * * * nvidia-smi --query-gpu=utilization.gpu,memory.used,temperature.gpu --format=csv,noheader >> /var/log/gpu-stats.log
```

## Next Steps

✅ GPU is now configured for Jellyfin and Immich!

**Recommended:**
1. Test transcoding with various file formats
2. Upload photos to Immich and verify ML inference works
3. Monitor GPU temperature and utilization
4. Consider unlocking NVENC stream limit
5. Set up automated monitoring

**Optional:**
- Configure Tdarr for batch transcoding using GPU
- Set up Plex (also supports NVENC)
- Use GPU for other workloads (AI, rendering)

## Reference

### Quick Command Reference

```bash
# Check GPU from host (Proxmox)
lspci | grep -i nvidia

# Check GPU from VM
nvidia-smi

# Test Docker GPU access
docker run --rm --gpus all nvidia/cuda:12.2.0-base-ubuntu22.04 nvidia-smi

# Monitor GPU real-time
watch -n 1 nvidia-smi

# Check Jellyfin GPU usage
docker exec jellyfin nvidia-smi

# Restart Jellyfin with GPU
cd ~/homelab/compose/media/frontend/jellyfin
docker compose down && docker compose up -d

# View GPU processes
nvidia-smi pmon

# GPU temperature
nvidia-smi --query-gpu=temperature.gpu --format=csv,noheader
```

### GTX 1070 Specifications

- **Architecture**: Pascal (GP104)
- **CUDA Cores**: 1920
- **Memory**: 8GB GDDR5
- **Memory Bandwidth**: 256 GB/s
- **TDP**: 150W
- **NVENC**: 6th generation (H.264, HEVC)
- **NVDEC**: 2nd generation
- **Concurrent Streams**: 2 (unlockable to unlimited)

---

**Your GTX 1070 is now accelerating your homelab! 🚀**
