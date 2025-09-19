package main

import (
	"fmt"
	"mvdan.cc/xurls/v2"
	"regexp"
)

// a helper function to combine and deduplicate string slices
func deduplicateAndCombine(slices ...[]string) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, slice := range slices {
		for _, item := range slice {
			if _, ok := seen[item]; !ok {
				seen[item] = struct{}{}
				result = append(result, item)
			}
		}
	}
	return result
}

func main() {
	text := `
		混合场景文本:
		1. 标准链接: https://google.com
		2. 无协议头链接: example.com and t.me/Emyzzy1
		3. 自定义App链接: bnc://user/profile/12345
		3. 自定义App链接: xxx://user/profile/12345
		4. 另一个未知协议: unknown-scheme://some/data/here?key=value
		5. 重复链接: https://google.com
		6. 重复无协议头链接: example.com and dshu.cuhu.com
	`

	// 第三步：合并与去重
	allLinks := ExtractWebLinks(text)

	fmt.Println("\n--- 最终合并去重后的结果 ---")
	for _, link := range allLinks {
		fmt.Println(link)
	}
}

func ExtractWebLinks(text string) []string {
	// 第一遍：使用 xurls.Relaxed() 提取标准和无协议头链接
	rxRelaxed := xurls.Relaxed()
	xurlsLinks := rxRelaxed.FindAllString(text, -1)
	fmt.Println("--- xurls 提取结果 ---")
	fmt.Println(xurlsLinks)

	// 第二遍：使用自定义正则提取所有带协议头的链接
	customRe := regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9+.-]*://\S+`)
	customLinks := customRe.FindAllString(text, -1)
	fmt.Println("\n--- 自定义正则提取结果 ---")
	fmt.Println(customLinks)

	// 第三步：合并与去重
	allLinks := deduplicateAndCombine(xurlsLinks, customLinks)

	fmt.Println("\n--- 最终合并去重后的结果 ---")
	return allLinks
}
