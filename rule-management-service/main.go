package main

import (
	"net/http"
	"rule-management-service/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

const rulesFile = "rules.json"

func main() {
	r := gin.Default()

	r.GET("/api/rules", func(c *gin.Context) {
		rules, err := utils.ReadRulesFromJSON(rulesFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read rules"})
			return
		}
		c.JSON(http.StatusOK, rules)
	})

	r.GET("/api/rules/:id", func(c *gin.Context) {
		ruleID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
			return
		}

		rules, err := utils.ReadRulesFromJSON(rulesFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read rules"})
			return
		}

		for _, rule := range rules {
			if rule.ID == ruleID {
				c.JSON(http.StatusOK, rule)
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Rule not found"})
	})

	// Add a new rule
	r.POST("/api/rules", func(c *gin.Context) {
		var newRule utils.Rule
		if err := c.ShouldBindJSON(&newRule); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		rules, err := utils.ReadRulesFromJSON(rulesFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load rules"})
			return
		}

		newRule.ID = utils.NextRuleID
		utils.NextRuleID++
		rules = append(rules, newRule)

		if err := utils.WriteRulesToJSON(rulesFile, rules); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save rules"})
			return
		}

		c.JSON(http.StatusCreated, newRule)
	})

	// Update an existing rule
	r.PUT("/api/rules/:id", func(c *gin.Context) {
		ruleID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
			return
		}

		var updatedRule utils.Rule
		if err := c.ShouldBindJSON(&updatedRule); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		rules, err := utils.ReadRulesFromJSON(rulesFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load rules"})
			return
		}

		updated := false
		for i, rule := range rules {
			if rule.ID == ruleID {
				updatedRule.ID = ruleID
				rules[i] = updatedRule
				updated = true
				break
			}
		}

		if !updated {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rule not found"})
			return
		}

		if err := utils.WriteRulesToJSON(rulesFile, rules); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save rules"})
			return
		}

		c.JSON(http.StatusOK, updatedRule)
	})

	// Delete a rule
	r.DELETE("/api/rules/:id", func(c *gin.Context) {
		ruleID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
			return
		}

		rules, err := utils.ReadRulesFromJSON(rulesFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load rules"})
			return
		}

		newRules := []utils.Rule{}
		deleted := false
		for _, rule := range rules {
			if rule.ID != ruleID {
				newRules = append(newRules, rule)
			} else {
				deleted = true
			}
		}

		if !deleted {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rule not found"})
			return
		}

		if err := utils.WriteRulesToJSON(rulesFile, newRules); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save rules"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Rule deleted"})
	})

	r.Run(":8004")

}
