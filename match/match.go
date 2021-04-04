package match

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "sync"
)

type item struct {
    code uint8
    x, y int
}

// P is pattern type for internal representation
type P map[int][]item

// ParseIncomingPattern returns P type to use it as match
func ParseIncomingPattern(src string) (P, error) {
    f, err := os.Open(src)
    if err != nil {
        return nil, err
    }

    defer f.Close()
    patternMap := make(P)
    scanner := bufio.NewScanner(f)
    curY := 0
    mapCounter := 0
    for scanner.Scan() {
        line := scanner.Text()
        tmpItems := []item{}
        for i := 0; i < len(line); i++ {
            if line[i] == ' ' {
                continue
            }
            tmpItems = append(tmpItems, item{code: line[i], x: i, y: curY})
        }
        if len(tmpItems) > 0 {
            patternMap[mapCounter] = tmpItems
            mapCounter++
        }
        curY++
    }
    return patternMap, nil
}

// FindByPattern find all occurrences in file
func FindByPattern(patternMap P, src string) (int, error) {
    tmpItems, err := findHeads(src, patternMap)
    if err != nil {
        return 0, err
    }
    if len(tmpItems) < 1 {
        return 0, nil
    }

    found := matchFull(tmpItems, patternMap, src)

    return found, nil
}

func findHeads(src string, patternMap P) ([]item, error) {
    hLine, err := getHeadLineFromPattern(patternMap)
    if err != nil {
        return nil, err
    }

    f, err := os.Open(src)
    if err != nil {
        return nil, err
    }

    defer f.Close()
    scanner := bufio.NewScanner(f)

    tmpItems := []item{}
    curY := 0
    for scanner.Scan() {
        line := scanner.Text()

        for i := 0; i < len(line); i++ {
            ok := findFirstHeadMatch(hLine[0].code, line, i)
            if !ok {
                continue
            }
            ok = findFullHeadMatch(hLine, line, i)
            if ok {
                tmpItems = append(tmpItems, item{x:i, y: curY})
            }
        }
        curY++
    }
    return tmpItems, nil
}

func findFullHeadMatch(head []item, line string, startFrom int) bool {
    proceed := false
    for _, h := range head {
        if proceed {
            startFrom += h.x
        }
        if startFrom >= len(line) {
            return false
        }
        if h.code != line[startFrom] {
            return false
        }
        proceed = true
    }
    return true
}

func findFirstHeadMatch(c uint8, line string, startFrom int) bool {
    if c == line[startFrom] {
        return true
    }
    return false
}

func getHeadLineFromPattern(patternMap P) ([]item, error) {
    v, ok := patternMap[0]
    if !ok {
        return nil, fmt.Errorf("%s", "Can't find items to search")
    }
    return v, nil
}

func matchFull(in []item, patternMap P, src string) int {
    if len(patternMap) == 1 {
        return len(in)
    }
    totalFound := 0
    wg := &sync.WaitGroup{}
    mu := &sync.Mutex{}
    for _, vIn := range in {
        zx := vIn.x - patternMap[0][0].x
        zy := vIn.y - patternMap[0][0].y
        rawPatternMap := make(P)
        curMapCur := 0

        for i := 1; i < len(patternMap); i++ {
            t := []item{}
            for _, v := range patternMap[i] {
                t = append(t, item{x: zx + v.x, y: zy + v.y, code: v.code})
            }
            rawPatternMap[curMapCur] = t
            curMapCur++
        }
        wg.Add(1)
        go func() {
            if ok := findInFile(rawPatternMap, src); ok {
                mu.Lock()
                totalFound++
                mu.Unlock()
            }
            wg.Done()
        }()
    }
    wg.Wait()
    return totalFound
}

func findInFile(m P, src string) bool {
    file, err := os.Open(src)
    if err != nil {
        log.Printf("Can't open file %q\n", src)
        return false
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    min := m[0][0].y
    max := m[len(m)-1][0].y

    sourceData := make(map[int]string)

    n := 0
    for scanner.Scan() {
        if n > max {
            break
        }
        if n < min {
            n++
            continue
        }
        for i := 0; i < len(m); i++ {
            if n == m[i][0].y {
                line := scanner.Text()
                sourceData[m[i][0].y] = line
                break
            }
        }
        n++
    }

    if len(m) != len(sourceData) {
        return false
    }
    for i := 0; i < len(m); i++ {
        for _, it := range m[i] {
            w := len(sourceData[it.y])
            if it.x >= w {
                return false
            }
            if it.x < 0 {
                return false
            }
            if it.code != sourceData[it.y][it.x] {
                return false
            }
        }
    }
    return true
}
