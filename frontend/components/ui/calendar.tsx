"use client"

import * as React from "react"
import { ChevronLeft, ChevronRight, ChevronDown } from "lucide-react"
import { DayPicker } from "react-day-picker"
import { ptBR } from "date-fns/locale"

import { cn } from "@/lib/utils"
import { buttonVariants } from "@/components/ui/button"

export type CalendarProps = React.ComponentProps<typeof DayPicker>

function Calendar({
  className,
  classNames,
  showOutsideDays = true,
  ...props
}: CalendarProps) {
  return (
    <DayPicker
      showOutsideDays={showOutsideDays}
      className={cn("p-4", className)}
      locale={ptBR}
      classNames={{
        root: "relative",
        months: "flex flex-col sm:flex-row gap-4",
        month: "flex flex-col gap-4",
        month_caption: "flex justify-center items-center h-10 px-10",
        caption_label: "text-sm font-medium hidden",
        nav: "absolute top-4 inset-x-4 h-10 flex items-center justify-between pointer-events-none z-10",
        button_previous: cn(
          buttonVariants({ variant: "outline" }),
          "h-7 w-7 bg-transparent p-0 opacity-50 hover:opacity-100 pointer-events-auto"
        ),
        button_next: cn(
          buttonVariants({ variant: "outline" }),
          "h-7 w-7 bg-transparent p-0 opacity-50 hover:opacity-100 pointer-events-auto"
        ),
        dropdowns: "flex items-center justify-center gap-2",
        dropdown: "appearance-none bg-transparent font-medium text-sm cursor-pointer focus:outline-none capitalize border rounded-md px-2 py-1 hover:bg-accent",
        dropdown_root: "relative inline-flex items-center",
        months_dropdown: "pr-6 capitalize",
        years_dropdown: "pr-6 max-h-48 overflow-y-auto",
        month_grid: "w-full border-collapse",
        weekdays: "flex",
        weekday:
          "text-muted-foreground rounded-md w-10 font-normal text-[0.8rem] text-center flex-1 uppercase",
        week: "flex w-full mt-2",
        day: "h-10 w-10 text-center text-sm p-0 relative flex-1 [&:has([aria-selected].day-range-end)]:rounded-r-md [&:has([aria-selected].day-outside)]:bg-accent/50 [&:has([aria-selected])]:bg-accent first:[&:has([aria-selected])]:rounded-l-md last:[&:has([aria-selected])]:rounded-r-md focus-within:relative focus-within:z-20",
        day_button: cn(
          buttonVariants({ variant: "ghost" }),
          "h-10 w-10 p-0 font-normal aria-selected:opacity-100 w-full"
        ),
        range_end: "day-range-end",
        selected:
          "bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground focus:bg-primary focus:text-primary-foreground",
        today: "bg-accent text-accent-foreground",
        outside:
          "day-outside text-muted-foreground opacity-50 aria-selected:bg-accent/50 aria-selected:text-muted-foreground aria-selected:opacity-30",
        disabled: "text-muted-foreground opacity-50",
        range_middle:
          "aria-selected:bg-accent aria-selected:text-accent-foreground",
        hidden: "invisible",
        chevron: "inline-block ml-1 h-4 w-4 opacity-50",
        ...classNames,
      }}
      components={{
        Chevron: ({ orientation }) => {
          if (orientation === "left") {
            return <ChevronLeft className="h-4 w-4" />
          }
          if (orientation === "right") {
            return <ChevronRight className="h-4 w-4" />
          }
          return <ChevronDown className="h-4 w-4 absolute right-0 top-1/2 -translate-y-1/2 pointer-events-none" />
        },
      }}
      {...props}
    />
  )
}
Calendar.displayName = "Calendar"

export { Calendar }
