// +build ignore

#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/udp.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>

#define ABTP_MAGIC 0x41425450 // "ABTP"
#define ABTP_PORT 9000
#define ABTP_MIN_LEN 8

// Minimal ABTP Header struct
struct abtp_hdr {
    __u32 magic;
    __u16 version;
    __u16 length;
};

SEC("xdp")
int abtp_xdp_prog(struct xdp_md *ctx) {
    void *data_end = (void *)(long)ctx->data_end;
    void *data = (void *)(long)ctx->data;

    // Parse Ethernet header
    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end) {
        return XDP_PASS;
    }

    if (eth->h_proto != bpf_htons(ETH_P_IP)) {
        return XDP_PASS;
    }

    // Parse IP header
    struct iphdr *ip = (void *)(eth + 1);
    if ((void *)(ip + 1) > data_end) {
        return XDP_PASS;
    }

    if (ip->protocol != 17) {
        return XDP_PASS;
    }

    // Parse UDP header
    struct udphdr *udp = (void *)(ip + 1);
    if ((void *)(udp + 1) > data_end) {
        return XDP_PASS;
    }

    if (udp->dest != bpf_htons(ABTP_PORT)) {
        return XDP_PASS;
    }

    // Parse ABTP header
    struct abtp_hdr *abtp = (void *)(udp + 1);
    if ((void *)(abtp + 1) > data_end) {
        // Frame too short
        return XDP_DROP;
    }

    // Validate Magic
    if (bpf_ntohl(abtp->magic) != ABTP_MAGIC) {
        // Invalid magic bytes
        return XDP_DROP;
    }

    // Validate Version (assume version 1 is supported)
    if (bpf_ntohs(abtp->version) != 1) {
        // Unsupported version
        return XDP_DROP;
    }

    // Validate Length (basic check against packet size)
    __u16 frame_len = bpf_ntohs(abtp->length);
    if (frame_len < ABTP_MIN_LEN) {
         // Malformed length
         return XDP_DROP;
    }

    // Valid ABTP frame, pass it to user-space / networking stack
    return XDP_PASS;
}

char _license[] SEC("license") = "GPL";
