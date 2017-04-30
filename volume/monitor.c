#include <stdio.h>
#include <errno.h>
#include <math.h>

#include <alsa/asoundlib.h>

#include "monitor.h"

static snd_ctl_t *open_ctl(int card) {
	snd_ctl_t *ctl;
	int err;

	// Find first available sound card if no specific index is set.
	if (card < 0) {
		snd_card_next(&card);
		if (card < 0) {
			fprintf(stderr, "No sound cards found\n");
			return NULL;
		}
	}

	// Canonical sound card name is `hw:<index>`.
	char name[16];
	sprintf(name, "hw:%d", card);

	// Open sound card with specific name.
	err = snd_ctl_open(&ctl, name, SND_CTL_READONLY);
	if (err < 0) {
		fprintf(stderr, "Cannot open sound card '%s'\n", name);
		return NULL;
	}

	err = snd_ctl_subscribe_events(ctl, 1);
	if (err < 0) {
		snd_ctl_close(ctl);
		fprintf(stderr, "Cannot open subscribe events to sound card '%s'\n", name);
		return NULL;
	}

	return ctl;
}

int volume() {
    snd_mixer_t *mixer;
    snd_mixer_elem_t *elem;
    snd_mixer_selem_id_t *id;

    snd_mixer_selem_id_alloca(&id);
    snd_mixer_selem_id_set_index(id, 0);
    snd_mixer_selem_id_set_name(id, "Master");

	if ((snd_mixer_open(&mixer, 0)) < 0) {
		fprintf(stderr, "Failed to open mixer\n");
		return -1;
	}

    if ((snd_mixer_attach(mixer, "default")) < 0) {
		fprintf(stderr, "Failed to attach mixer\n");
        snd_mixer_close(mixer);
        return -1;
    }

    if ((snd_mixer_selem_register(mixer, NULL, NULL)) < 0) {
		fprintf(stderr, "Failed to register mixer element\n");
        snd_mixer_close(mixer);
        return -1;
    }

    if ((snd_mixer_load(mixer)) < 0) {
		fprintf(stderr, "Failed to load mixer element\n");
        snd_mixer_close(mixer);
        return -1;
    }

    elem = snd_mixer_find_selem(mixer, id);
    if (!elem) {
		fprintf(stderr, "Failed to find mixer element\n");
        snd_mixer_close(mixer);
        return -1;
    }

    long volume, volume_min, volume_max;
    snd_mixer_selem_get_playback_volume_range(elem, &volume_min, &volume_max);

	if (snd_mixer_selem_get_playback_volume(elem, 0, &volume) < 0) {
		snd_mixer_close(mixer);
		return -1;
	}

	volume -= volume_min;
	volume_max -= volume_min;

	snd_mixer_close(mixer);

	return ceil((double) volume / (double) volume_max * 100);
}

void wait() {
	snd_ctl_t *ctl = open_ctl(-1);
	if (ctl == NULL) {
		errno = 1;
		return;
	}

	int err;
	struct pollfd fd;

	snd_ctl_poll_descriptors(ctl, &fd, 1);

	err = poll(&fd, 1, -1);
	if (err <= 0) {
		snd_ctl_close(ctl);
		errno = err;
		return;
	}

	unsigned short revents;
	snd_ctl_poll_descriptors_revents(ctl, &fd, 1, &revents);
	snd_ctl_close(ctl);
}
