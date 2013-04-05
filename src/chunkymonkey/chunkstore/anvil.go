package chunkstore

import (
    "fmt"
    "log"
    "os"
    "path"

    . "chunkymonkey/types"
)

const (
    regionFileEdge       = 32
    regionFileEdgeShift  = 5
    regionFileSectorSize = 4096
    // 5 is the size of chunkDataHeader in bytes.
    chunkDataHeaderSize = 5
    chunkDataGuessSize  = 8192

    chunkCompressionGzip = 1
    chunkCompressionZlib = 2
)

type chunkStoreAnvil struct {
    regionPath  string
    regionFiles map[uint64]*regionFile
}

func newChunkStoreAnvil(worldPath string, dimension DimensionId) (s *chunkStoreAnvil, err error) {
    s = &chunkStoreAnvil{
        regionFiles: make(map[uint64]*regionFile),
    }

    if dimension == DimensionNormal {
        s.regionPath = path.Join(worldPath, "region")
    } else {
        s.regionPath = path.Join(worldPath, fmt.Sprintf("DIM%d", dimension), "region")
    }

    if err = os.MkdirAll(s.regionPath, 0777); err != nil {
        return nil, err
    }

    return
}

func (s *chunkStoreAnvil) regionFile(chunkLoc ChunkXz) (rf *regionFile, err error) {
    regionLoc := regionLocForChunkXz(chunkLoc)

    rf, ok := s.regionFiles[regionLoc.regionKey()]
    if ok {
        return rf, nil
    }

    // TODO limit number of regionFile objs to a maximum number of
    // most-frequently-used regions. Close regionFile objects when no
    // longer needed.
    filePath := regionLoc.regionFilePath(s.regionPath)
    log.Printf("PATH: %s", filePath)
    rf, err = newRegionFile(filePath)
    if err != nil {
        if os.IsNotExist(err) {
            err = NoSuchChunkError(false)
        }
        return
    }
    s.regionFiles[regionLoc.regionKey()] = rf

    return rf, nil
}

func (s *chunkStoreAnvil) ReadChunk(chunkLoc ChunkXz) (reader IChunkReader, err error) {
    rf, err := s.regionFile(chunkLoc)
    if err != nil {
        return
    }

    chunkReader, err := rf.ReadChunkData(chunkLoc)
    if chunkReader != nil {
        reader = chunkReader
    }

    return
}

func (s *chunkStoreAnvil) SupportsWrite() bool {
    return true
}

func (s *chunkStoreAnvil) Writer() IChunkWriter {
    return newNbtChunkWriter()
}

func (s *chunkStoreAnvil) WriteChunk(writer IChunkWriter) error {
    nbtWriter, ok := writer.(*nbtChunkWriter)
    if !ok {
        return fmt.Errorf("%T is incorrect IChunkWriter implementation for %T", writer, s)
    }

    rf, err := s.regionFile(writer.ChunkLoc())
    if err != nil {
        return err
    }

    return rf.WriteChunkData(nbtWriter)
}
