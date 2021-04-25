CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE tokens (
    id UUID NOT NULL,
    access_token TEXT NOT NULL,
    token_type TEXT NOT NULL,
    scope TEXT NOT NULL,
    expires_in INT NOT NULL,
    refresh_token TEXT NOT NULL,

    PRIMARY KEY(id)
);

CREATE TABLE users (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    display_name TEXT NOT NULL,
    email TEXT NOT NULL,
    spotify_url TEXT NOT NULL,
    spotify_image_url TEXT NOT NULL,
    spotify_id TEXT NOT NULL,
    token TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    PRIMARY KEY(id)
);

CREATE TABLE playlists (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    spotify_url TEXT NOT NULL,
    spotify_id TEXT NOT NULL,
    total_tracks INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    user_id UUID NOT NULL,

    PRIMARY KEY(id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
        REFERENCES users (id) MATCH SIMPLE
        ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE tracks (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    popularity INT NOT NULL,
    duration INT NOT NULL,
    added_at TEXT NOT NULL,
    spotify_uri TEXT NOT NULL,
    spotify_url TEXT NOT NULL,
    spotify_id TEXT NOT NULL,
    artists_json JSON NOT NULL,
    album_json JSON NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    PRIMARY KEY(id)
);

CREATE TABLE playlists_tracks (
    playlist_id UUID NOT NULL,
    track_id UUID NOT NULL,

    PRIMARY KEY(playlist_id, track_id),
    CONSTRAINT fk_playlist_id FOREIGN KEY (playlist_id)
        REFERENCES playlists (id) MATCH SIMPLE
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_track_id FOREIGN KEY (track_id)
        REFERENCES tracks (id) MATCH SIMPLE
        ON UPDATE CASCADE ON DELETE CASCADE
);
