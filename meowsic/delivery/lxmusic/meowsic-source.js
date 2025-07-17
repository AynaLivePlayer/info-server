/*!
 * @name Meowsic
 * @description Meowsic Source
 * @version v0.0.1
 * @author Aynakeya
 * Code are adapt from https://github.com/Huibq/keep-alive/blob/master/render_api.js
 */
const DEV_ENABLE = false
const API_URL = 'http://127.0.0.1:9090'
const MUSIC_QUALITY = {
    kw: ['128k', '320k'],
    kg: ['128k', '320k'],
    tx: ['128k', '320k'],
    wy: ['128k', '320k'],
}
const MUSIC_SOURCE = Object.keys(MUSIC_QUALITY)
const { EVENT_NAMES, request, on, send, utils, env, version } = globalThis.lx
const httpFetch = (url, options = { method: 'GET' }) => {
    return new Promise((resolve, reject) => {
        request(url, options, (err, resp) => {
            if (err) return reject(err)
            resolve(resp)
        })
    })
}
const handleGetMusicUrl = async (source, musicInfo, quality) => {
    const songId = musicInfo.hash ?? musicInfo.songmid

    const request = await httpFetch(`${API_URL}/api/meowsic/lxmusic/url?provider=${source}&id=${songId}&quality=${quality}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'User-Agent': `${env ? `lx-music-${env}/${version}` : `lx-usic-request/${version}`}`,
        },
    })
    const { body } = request
    if (!body || isNaN(Number(body.code))) throw new Error('unknow error')
    if (DEV_ENABLE) {
        console.log(body);
    }
    switch (body.code) {
        case 0:
            return body.data[0].url
        default:
            throw new Error(body.msg ?? 'unknow error')
    }
}
const musicSources = {}
MUSIC_SOURCE.forEach(item => {
    musicSources[item] = {
        name: item,
        type: 'music',
        actions: ['musicUrl'],
        qualitys: MUSIC_QUALITY[item],
    }
})
on(EVENT_NAMES.request, ({ action, source, info }) => {
    switch (action) {
        case 'musicUrl':
            if (env != 'mobile') {
                console.group(`Handle Action(musicUrl)`)
                console.log('source', source)
                console.log('quality', info.type)
                console.log('musicInfo', info.musicInfo)
                console.groupEnd()
            } else {
                console.log(`Handle Action(musicUrl)`)
                console.log('source', source)
                console.log('quality', info.type)
                console.log('musicInfo', info.musicInfo)
            }
            return handleGetMusicUrl(source, info.musicInfo, info.type)
                .then(data => Promise.resolve(data))
                .catch(err => Promise.reject(err))
        default:
            console.error(`action(${action}) not support`)
            return Promise.reject('action not support')
    }
})
send(EVENT_NAMES.inited, { status: true, openDevTools: DEV_ENABLE, sources: musicSources })