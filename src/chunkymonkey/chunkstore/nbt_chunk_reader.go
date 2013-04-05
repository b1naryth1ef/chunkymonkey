package chunkstore

import (
    "io"
    "log"

    "chunkymonkey/gamerules"
    . "chunkymonkey/types"
    "nbt"
)

// Returned to chunks to pull their data from.
type nbtChunkReader struct {
    chunkTag nbt.ITag
}

// Combine arrays
func combineArray(arr1, arr2 []byte) []byte {
    n := len(arr1) + len(arr2)
    if n > cap(arr1) {
        newSlice := make([]byte, (n+1)*2)
        copy(newSlice, arr1)
        arr1 = newSlice
    }
    copy(arr1[len(arr1):n], arr2)
    return arr1
}

// Load a chunk from its NBT representation
func newNbtChunkReader(reader io.Reader) (r *nbtChunkReader, err error) {
    chunkTag, err := nbt.Read(reader)
    if err != nil {
        return
    }

    r = &nbtChunkReader{
        chunkTag: chunkTag,
    }

    return
}

func (r *nbtChunkReader) ChunkLoc() ChunkXz { //@FIXME
    _, ok := r.chunkTag.Lookup("Level/zPos").(*nbt.Int)
    if !ok {
        log.Panicln("Invalid Level Format!")
    }
    return ChunkXz{
        X:  ChunkCoord(r.chunkTag.Lookup("Level/xPos").(*nbt.Int).Value),
        Z:  ChunkCoord(r.chunkTag.Lookup("Level/zPos").(*nbt.Int).Value),
    }
}

func (r *nbtChunkReader) Sections() (comps []nbt.Compound) {
    sectionsTag, ok := r.chunkTag.Lookup("Level/Sections").(*nbt.List)
    if !ok {
        log.Printf("Failed to load Sections %s", sectionsTag)
        return
    }

    comps = make([]nbt.Compound, len(sectionsTag.Value))
    for _, value := range sectionsTag.Value {
        comp, ok := value.(nbt.Compound)
        if !ok {
            log.Printf("Found non-compound in sections list: %T", value)
            continue
        }

        comps = append(comps, comp)
    }

    return comps
}

func (r *nbtChunkReader) Blocks() []byte {
    secs := r.Sections()
    res := make([]byte, 4096)

    for _, value := range secs {
        v := value.Lookup("Blocks").(*nbt.ByteArray).Value
        log.Printf("%s", v)
        res = combineArray(res, v)
    }

    return res
}

func (r *nbtChunkReader) BlockData() []byte {
    secs := r.Sections()
    res := make([]byte, 4096*len(secs))

    for _, value := range secs {
        copy(res, value.Lookup("BlockData").(*nbt.ByteArray).Value)
    }

    return res
}

func (r *nbtChunkReader) BlockLight() []byte {
    secs := r.Sections()
    res := make([]byte, 4096*len(secs))

    for _, value := range secs {
        copy(res, value.Lookup("BlockLight").(*nbt.ByteArray).Value)
    }

    return res
}

func (r *nbtChunkReader) SkyLight() []byte {
    secs := r.Sections()
    res := make([]byte, 4096*len(secs))

    for _, value := range secs {
        copy(res, value.Lookup("SkyLight").(*nbt.ByteArray).Value)
    }

    return res
}

func (r *nbtChunkReader) HeightMap() []int {
    return r.chunkTag.Lookup("Level/HeightMap").(*nbt.IntArray).Value
}

func (r *nbtChunkReader) Entities() (entities []gamerules.INonPlayerEntity) {
    entityListTag, ok := r.chunkTag.Lookup("Level/Entities").(*nbt.List)
    if !ok {
        return
    }

    entities = make([]gamerules.INonPlayerEntity, 0, len(entityListTag.Value))
    for _, tag := range entityListTag.Value {
        compound, ok := tag.(nbt.Compound)
        if !ok {
            log.Printf("Found non-compound in entities list: %T", tag)
            continue
        }

        entityObjectId, ok := compound.Lookup("id").(*nbt.String)
        if !ok {
            log.Printf("Missing or bad entity type ID in NBT: %s", entityObjectId)
            continue
        }

        entity := gamerules.NewEntityByTypeName(entityObjectId.Value)
        if entity == nil {
            log.Printf("Found unhandled entity type: %s", entityObjectId.Value)
            continue
        }

        err := entity.UnmarshalNbt(compound)
        if err != nil {
            log.Printf("Error unmarshalling entity NBT: %s", err)
            continue
        }

        entities = append(entities, entity)
    }

    return
}

func (r *nbtChunkReader) TileEntities() (tileEntities []gamerules.ITileEntity) {
    entityListTag, ok := r.chunkTag.Lookup("Level/TileEntities").(*nbt.List)
    if !ok {
        return
    }

    tileEntities = make([]gamerules.ITileEntity, 0, len(entityListTag.Value))
    for _, tag := range entityListTag.Value {
        compound, ok := tag.(nbt.Compound)
        if !ok {
            log.Printf("Found non-compound in tile entities list: %T", tag)
            continue
        }

        entityObjectId, ok := compound.Lookup("id").(*nbt.String)
        if !ok {
            log.Printf("Missing or bad tile entity type ID in NBT: %s", entityObjectId)
            continue
        }

        entity := gamerules.NewTileEntityByTypeName(entityObjectId.Value)
        if entity == nil {
            log.Printf("Found unhandled tile entity type: %s", entityObjectId.Value)
            continue
        }

        if err := entity.UnmarshalNbt(compound); err != nil {
            log.Printf("%T.UnmarshalNbt failed for %s: %s", entity, compound, err)
            continue
        }

        tileEntities = append(tileEntities, entity)
    }

    return
}

func (r *nbtChunkReader) RootTag() nbt.ITag {
    return r.chunkTag
}
