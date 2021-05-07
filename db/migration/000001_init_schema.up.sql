CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE tokens (
    id UUID NOT NULL,
    access_token TEXT NOT NULL,
    token_type TEXT NOT NULL,
    scope TEXT NOT NULL,
    expires_in INT NOT NULL,
    refresh_token TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    PRIMARY KEY(id)
);

CREATE TABLE users (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    display_name TEXT NOT NULL,
    email TEXT NOT NULL,
    spotify_url TEXT NOT NULL,
    spotify_image_url TEXT NOT NULL,
    spotify_id TEXT NOT NULL,
    auth_id TEXT NOT NULL,
    token TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    PRIMARY KEY(id)
);

-- @TODO: dummy data, remove
INSERT INTO users (
    id,
    display_name,
    email,
    spotify_url,
    spotify_image_url,
    spotify_id,
    auth_id,
    token
) VALUES (
    '12345678-0621-46b7-881a-69c6959d65e1',
    'Joe Smith',
    'joe@email.com',
    'foo',
    'foo',
    'foo',
    'foo',
    'foo'
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

-- @TODO: dummy data, remove
INSERT INTO playlists (
    id,
    name,
    spotify_url,
    spotify_id,
    total_tracks,
    user_id
) VALUES (
    '12345678-8a29-4113-9565-349e4584e640',
    'Joe Playlist',
    'foo',
    'foo',
    1,
    '12345678-0621-46b7-881a-69c6959d65e1'
);

CREATE TABLE tracks (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    popularity INT NOT NULL,
    duration INT NOT NULL,
    spotify_uri TEXT NOT NULL,
    spotify_url TEXT NOT NULL,
    spotify_id TEXT NOT NULL,
    artists_json JSONB NOT NULL,
    album_json JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    PRIMARY KEY(id)
);

-- @TODO: dummy data, remove
INSERT INTO tracks (
    id,
    name,
    popularity,
    duration,
    spotify_uri,
    spotify_url,
    spotify_id,
    artists_json,
    album_json
) VALUES (
    '35ab5489-57a3-46ce-b65f-506c86ef94f8',
    'Undercurrents',
    24,
    285210,
    'spotify:track:7wNEQEkg7Nh5NPDgti2DF9',
    'https://open.spotify.com/track/7wNEQEkg7Nh5NPDgti2DF9',
    '7wNEQEkg7Nh5NPDgti2DF9',
    '[{"name":"Frost","spotify_url":"https://open.spotify.com/artist/4cr1vZsdjcY434Aqc3fDBt"}]'::JSON,
    '{"name":"Anjunadeep Explorations 12","spotify_url":"https://open.spotify.com/album/1lpShXCM4PHvjsLBQpzFJY","spotify_image_url":"https://i.scdn.co/image/ab67616d0000b273d3bc3e38ad7b813b4d9e4cb7"}'::JSON
);

CREATE TABLE playlists_tracks (
    playlist_id UUID NOT NULL,
    track_id UUID NOT NULL,
    added_at TEXT NOT NULL,

    PRIMARY KEY(playlist_id, track_id),
    CONSTRAINT fk_playlist_id FOREIGN KEY (playlist_id)
        REFERENCES playlists (id) MATCH SIMPLE
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_track_id FOREIGN KEY (track_id)
        REFERENCES tracks (id) MATCH SIMPLE
        ON UPDATE CASCADE ON DELETE CASCADE
);

-- @TODO: dummy data, remove
INSERT INTO playlists_tracks (
    playlist_id,
    track_id,
    added_at
) VALUES (
    '12345678-8a29-4113-9565-349e4584e640',
    '35ab5489-57a3-46ce-b65f-506c86ef94f8',
    '2021-04-27T17:32:31Z'
);

CREATE TABLE friendships (
    user_id UUID NOT NULL,
    friend_id UUID NOT NULL,
    pending BOOLEAN NOT NULL DEFAULT true,

    PRIMARY KEY(user_id, friend_id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
        REFERENCES users (id) MATCH SIMPLE
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_friend_id FOREIGN KEY (friend_id)
        REFERENCES users (id) MATCH SIMPLE
        ON UPDATE CASCADE ON DELETE CASCADE
)
