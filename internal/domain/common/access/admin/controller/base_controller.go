package controller

// BaseController maps to BaseController.java
// @MappedFrom BaseController
type BaseController struct {
}

// @MappedFrom getPageSize(@Nullable Integer size)
func (c *BaseController) GetPageSize() {
}

// @MappedFrom queryBetweenDate(DateRange dateRange, DivideBy divideBy, Function3<DateRange, Boolean, Boolean, Mono<Long>> function, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages)
func (c *BaseController) QueryBetweenDate() {
}

// @MappedFrom queryBetweenDate(DateRange dateRange, DivideBy divideBy, Function<DateRange, Mono<Long>> function)
func (c *BaseController) QueryBetweenDateFunc() {
}

// @MappedFrom checkAndQueryBetweenDate(DateRange dateRange, DivideBy divideBy, Function3<DateRange, Boolean, Boolean, Mono<Long>> function, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages)
func (c *BaseController) CheckAndQueryBetweenDate() {
}

// @MappedFrom checkAndQueryBetweenDate(DateRange dateRange, DivideBy divideBy, Function<DateRange, Mono<Long>> function)
func (c *BaseController) CheckAndQueryBetweenDateFunc() {
}
